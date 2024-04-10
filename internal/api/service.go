package api

import (
	"golang-api/internal/db"
	"regexp"
)

type Service interface {
	Register(teacherEmail string, studentEmails []string) error
	CommonStudents(teacherEmails []string) ([]string, error)
	Suspend(studentEmail string) error
	RetrieveForNotifications(teacherEmail string, notification string) ([]string, error)
}

type ServiceImpl struct {
	repo db.Repository
}

func NewService(repo db.Repository) Service {
	return &ServiceImpl{repo: repo}
}

func (s *ServiceImpl) Register(teacherEmail string, studentEmails []string) error {
	// Lookup the teacher by email, or create a new one if it doesn't exist
	teacher, err := s.repo.FindOrCreateTeacherByEmail(teacherEmail)
	if err != nil {
		return err
	}

	// Lookup the students by email, or create new ones if they don't exist
	students := make([]db.Student, 0, len(studentEmails))
	for _, studentEmail := range studentEmails {
		student, err := s.repo.FindOrCreateStudentByEmail(studentEmail)
		if err != nil {
			return err
		}
		students = append(students, student)
	}

	// Associate the students with the teacher
	if err := s.repo.AssociateTeacherWithStudents(teacher, students); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) CommonStudents(teacherEmails []string) ([]string, error) {
	// Lookup the teachers by email
	teachers := make([]db.Teacher, 0, len(teacherEmails))
	for _, teacherEmail := range teacherEmails {
		teacher, err := s.repo.FindTeacherByEmail(teacherEmail)
		if err != nil {
			return nil, err
		}
		teachers = append(teachers, teacher)
	}

	// Find the common students for the teachers
	students, err := s.repo.FindCommonStudentsForTeachers(teachers)
	if err != nil {
		return nil, err
	}

	// Extract the student emails
	studentEmails := make([]string, 0, len(students))
	for _, student := range students {
		studentEmails = append(studentEmails, student.Email)
	}

	return studentEmails, nil
}

func (s *ServiceImpl) Suspend(studentEmail string) error {
	// Lookup the student by email
	student, err := s.repo.FindStudentByEmail(studentEmail)
	if err != nil {
		return err
	}

	// Suspend the student
	student.IsSuspended = true

	// Update the student
	if err := s.repo.UpdateStudent(student); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) RetrieveForNotifications(teacherEmail string, notification string) ([]string, error) {
	// Extract mentioned student emails from the notification
	mentionedStudentEmails := extractMentionedStudents(notification)

	// Find or create mentioned students
	mentionedStudents := make([]db.Student, 0, len(mentionedStudentEmails))
	for _, studentEmail := range mentionedStudentEmails {
		student, err := s.repo.FindOrCreateStudentByEmail(studentEmail)
		if err != nil {
			return nil, err
		}
		mentionedStudents = append(mentionedStudents, student)
	}

	// Lookup the teacher by email, or create a new one if it doesn't exist
	teacher, err := s.repo.FindOrCreateTeacherByEmail(teacherEmail)
	if err != nil {
		return nil, err
	}

	// Find students who are registered to this teacher
	teachers := []db.Teacher{teacher}
	registeredStudents, err := s.repo.FindCommonStudentsForTeachers(teachers)
	if err != nil {
		return nil, err
	}

	// Combine the mentioned students and registered students, removing duplicates
	students := make(map[string]db.Student)
	for _, student := range mentionedStudents {
		students[student.Email] = student
	}
	for _, student := range registeredStudents {
		students[student.Email] = student
	}

	// Remove suspended students
	for studentEmail, student := range students {
		if student.IsSuspended {
			delete(students, studentEmail)
		}
	}

	// Extract the student emails
	studentEmails := make([]string, 0, len(students))
	for studentEmail := range students {
		studentEmails = append(studentEmails, studentEmail)
	}

	return studentEmails, nil
}

// extractMentionedStudents finds all mentioned emails in the notification text
func extractMentionedStudents(notification string) []string {
	emailRegex := regexp.MustCompile(`@\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`)
	matches := emailRegex.FindAllString(notification, -1)

	for i, match := range matches {
		matches[i] = match[1:] // Remove the leading '@'
	}

	return matches
}

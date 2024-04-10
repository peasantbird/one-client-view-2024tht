package api

import "golang-api/internal/db"

type Service interface {
	Register(teacher string, students []string) error
}

type ServiceImpl struct {
	repo db.Repository
}

func NewService(repo db.Repository) Service {
	return &ServiceImpl{repo: repo}
}

func (s *ServiceImpl) Register(teacherEmail string, studentEmails []string) error {
	// Lookup the teacher by email, or create a new one if it doesn't exist
	teacher, err := s.repo.FindOrCreateTeacher(teacherEmail)
	if err != nil {
		return err
	}

	// Lookup the students by email, or create new ones if they don't exist
	students := make([]db.Student, 0, len(studentEmails))
	for _, studentEmail := range studentEmails {
		student, err := s.repo.FindOrCreateStudent(studentEmail)
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

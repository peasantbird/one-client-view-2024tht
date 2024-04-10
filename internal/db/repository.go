package db

type Repository interface {
	FindTeacherByEmail(email string) (Teacher, error)
	FindStudentByEmail(email string) (Student, error)
	FindOrCreateTeacherByEmail(email string) (Teacher, error)
	FindOrCreateStudentByEmail(email string) (Student, error)
	UpdateStudent(student Student) error
	AssociateTeacherWithStudents(teacher Teacher, students []Student) error
	FindCommonStudentsForTeachers(teachers []Teacher) ([]Student, error)
}

package db

type Repository interface {
	FindTeacherByEmail(email string) (Teacher, error)
	FindOrCreateTeacherByEmail(email string) (Teacher, error)
	FindOrCreateStudentByEmail(email string) (Student, error)
	AssociateTeacherWithStudents(teacher Teacher, students []Student) error
	FindCommonStudentsForTeachers(teachers []Teacher) ([]Student, error)
}

package db

type Repository interface {
	FindOrCreateTeacher(email string) (Teacher, error)
	FindOrCreateStudent(email string) (Student, error)
	AssociateTeacherWithStudents(teacher Teacher, students []Student) error
}

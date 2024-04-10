package db

import (
	"gorm.io/gorm"
)

type GormPostgresRepository struct {
	db *gorm.DB
}

func NewGormPostgresRepository(db *gorm.DB) *GormPostgresRepository {
	if db == nil {
		panic("Missing database connection")
	}

	return &GormPostgresRepository{db: db}
}

func (r *GormPostgresRepository) FindTeacherByEmail(email string) (Teacher, error) {
	var teacher Teacher
	if err := r.db.Where(Teacher{Email: email}).First(&teacher).Error; err != nil {
		return Teacher{}, err
	}

	return teacher, nil
}

func (r *GormPostgresRepository) FindOrCreateTeacherByEmail(email string) (Teacher, error) {
	var teacher Teacher
	if err := r.db.Where(Teacher{Email: email}).FirstOrCreate(&teacher).Error; err != nil {
		return Teacher{}, err
	}

	return teacher, nil
}

func (r *GormPostgresRepository) FindOrCreateStudentByEmail(email string) (Student, error) {
	var student Student
	if err := r.db.Where(Student{Email: email}).FirstOrCreate(&student).Error; err != nil {
		return Student{}, err
	}

	return student, nil
}

func (r *GormPostgresRepository) AssociateTeacherWithStudents(teacher Teacher, students []Student) error {
	if err := r.db.Model(&teacher).Association("Students").Append(&students); err != nil {
		return err
	}

	return nil
}

func (r *GormPostgresRepository) FindCommonStudentsForTeachers(teachers []Teacher) ([]Student, error) {
	teacherIDs := make([]uint, 0, len(teachers))
	for _, teacher := range teachers {
		teacherIDs = append(teacherIDs, teacher.ID)
	}

	var students []Student
	err := r.db.Model(&Student{}).Joins("JOIN teacher_students on teacher_students.student_id = students.id").
		Joins("JOIN teachers on teachers.id = teacher_students.teacher_id").
		Where("teachers.id IN ?", teacherIDs).
		Group("students.id").
		Having("COUNT(DISTINCT teachers.id) = ?", len(teachers)).
		Find(&students).Error

	return students, err
}

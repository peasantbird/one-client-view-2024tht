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

func (r *GormPostgresRepository) FindOrCreateTeacher(email string) (Teacher, error) {
	var teacher Teacher
	if err := r.db.Where(Teacher{Email: email}).FirstOrCreate(&teacher).Error; err != nil {
		return Teacher{}, err
	}

	return teacher, nil
}

func (r *GormPostgresRepository) FindOrCreateStudent(email string) (Student, error) {
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

package db

import (
	"gorm.io/gorm"
)

type Teacher struct {
	gorm.Model
	Email    string     `gorm:"unique;not null"`
	Students []*Student `gorm:"many2many:teacher_students;"`
}

type Student struct {
	gorm.Model
	Email       string     `gorm:"unique;not null"`
	IsSuspended bool       `gorm:"default:false;not null"`
	Teachers    []*Teacher `gorm:"many2many:teacher_students;"`
}

type TeacherStudent struct {
	TeacherID uint
	StudentID uint
}

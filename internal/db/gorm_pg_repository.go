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

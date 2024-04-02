package db

import (
	"golang-api/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(c *config.Config) (*gorm.DB, error) {
	dsn := c.DB.DSN

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&Teacher{}, &Student{})

	return db, nil
}

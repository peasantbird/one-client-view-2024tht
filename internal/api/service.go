package api

import "golang-api/internal/db"

type Service interface{}

type ServiceImpl struct {
	repo db.Repository
}

func NewService(repo db.Repository) Service {
	return &ServiceImpl{repo: repo}
}

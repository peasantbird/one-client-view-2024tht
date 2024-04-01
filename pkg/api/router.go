package api

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Router(h *Handler) *mux.Router {
	router := mux.NewRouter()

	apiRoute := router.PathPrefix("/api").Subrouter()
	apiRoute.HandleFunc("/register", h.Register).Methods("POST")

	return router
}

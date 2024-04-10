package api

import (
	"github.com/gorilla/mux"
)

func Router(h *Handler) *mux.Router {
	router := mux.NewRouter()

	apiRoute := router.PathPrefix("/api").Subrouter()
	apiRoute.HandleFunc("/register", h.Register).Methods("POST")
	// apiRoute.HandleFunc("/commonstudents", h.CommonStudents).Methods("GET")
	// apiRoute.HandleFunc("/suspend", h.Suspend).Methods("POST")
	// apiRoute.HandleFunc("/retrievefornotifications", h.RetrieveForNotifications).Methods("POST")

	return router
}

package api

import (
	"encoding/json"
	"golang-api/pkg/db"
	"net/http"
)

// RegisterRequest defines the expected format of the request body
type RegisterRequest struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the RegisterRequest struct
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lookup the teacher by email, or create a new one if it doesn't exist
	var teacher db.Teacher
	if err := h.DB.Where(db.Teacher{Email: req.Teacher}).FirstOrCreate(&teacher).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Process each student email
	for _, studentEmail := range req.Students {
		var student db.Student
		// Find or create the student record
		if err := h.DB.Where(db.Student{Email: studentEmail}).FirstOrCreate(&student).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Associate the student with the teacher
		// It will automatically handle the many-to-many relationship
		if err := h.DB.Model(&teacher).Association("Students").Append(&student); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Respond with HTTP 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}

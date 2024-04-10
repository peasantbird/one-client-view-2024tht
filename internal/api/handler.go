package api

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the request struct
	var req struct {
		Teacher  string   `json:"teacher"`
		Students []string `json:"students"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service method to register students to a teacher
	if err := h.service.Register(req.Teacher, req.Students); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with HTTP 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CommonStudents(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters
	teacherEmails := r.URL.Query()["teacher"]
	if len(teacherEmails) == 0 {
		http.Error(w, "No teacher query parameter provided", http.StatusBadRequest)
		return
	}

	// Call the service method to find common students for the teachers
	studentEmails, err := h.service.CommonStudents(teacherEmails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the response object
	type CommonStudentsResponse struct {
		Students []string `json:"students"`
	}
	response := CommonStudentsResponse{
		Students: studentEmails,
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) Suspend(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the request struct
	var req struct {
		Student string `json:"student"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service method to suspend a student
	if err := h.service.Suspend(req.Student); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with HTTP 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}

// type NotificationRequest struct {
// 	Teacher      string `json:"teacher"`
// 	Notification string `json:"notification"`
// }

// type NotificationResponse struct {
// 	Recipients []string `json:"recipients"`
// }

// func (h *Handler) RetrieveForNotifications(w http.ResponseWriter, r *http.Request) {
// 	// Decode the JSON body into the NotificationRequest struct
// 	var req NotificationRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	// Extract mentioned students from the notification text
// 	mentioned := extractMentionedStudents(req.Notification)

// 	// Get the list of students who are not suspended and are either registered to the teacher or mentioned in the notification
// 	var students []db.Student
// 	h.DB.Model(&db.Student{}).
// 		Select("students.email").
// 		Joins("LEFT JOIN teacher_students ON teacher_students.student_id = students.id").
// 		Joins("LEFT JOIN teachers ON teachers.id = teacher_students.teacher_id").
// 		Where("students.is_suspended = ?", false).
// 		Where("teachers.email = ? OR students.email IN ?", req.Teacher, mentioned).
// 		Group("students.email").
// 		Find(&students)

// 	// Prepare the list of student emails
// 	studentEmails := make([]string, len(students))
// 	for i, student := range students {
// 		studentEmails[i] = student.Email
// 	}

// 	// Create the response object
// 	response := NotificationResponse{
// 		Recipients: studentEmails,
// 	}

// 	// Write the response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(response); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// // extractMentionedStudents finds all mentioned emails in the notification text
// func extractMentionedStudents(notification string) []string {
// 	emailRegex := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
// 	return emailRegex.FindAllString(notification, -1)
// }

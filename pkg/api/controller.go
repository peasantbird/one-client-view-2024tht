package api

import (
	"encoding/json"
	"golang-api/pkg/db"
	"net/http"
	"regexp"
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

type CommonStudentsResponse struct {
	Students []string `json:"students"`
}

func (h *Handler) CommonStudents(w http.ResponseWriter, r *http.Request) {
	// Extract query parameters
	teacherEmails := r.URL.Query()["teacher"]
	if len(teacherEmails) == 0 {
		http.Error(w, "No teacher query parameter provided", http.StatusBadRequest)
		return
	}

	// Find common students for the given teachers
	var students []db.Student
	h.DB.Model(&db.Student{}).Joins("JOIN teacher_students on teacher_students.student_id = students.id").
		Joins("JOIN teachers on teachers.id = teacher_students.teacher_id").
		Where("teachers.email IN ?", teacherEmails).
		Group("students.id").
		Having("COUNT(DISTINCT teachers.id) = ?", len(teacherEmails)).
		Find(&students)

	// Prepare the list of student emails
	studentEmails := make([]string, 0, len(students))
	for _, student := range students {
		studentEmails = append(studentEmails, student.Email)
	}

	// Create the response object
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
	// Decode the JSON body into a map
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lookup the student by email
	var student db.Student
	if err := h.DB.Where(db.Student{Email: req["student"]}).First(&student).Error; err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}

	// Update the student's suspended status
	if err := h.DB.Model(&student).Update("IsSuspended", true).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with HTTP 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}

// NotificationRequest defines the expected format of the request body
type NotificationRequest struct {
	Teacher      string `json:"teacher"`
	Notification string `json:"notification"`
}

// NotificationResponse defines the format of the response body
type NotificationResponse struct {
	Recipients []string `json:"recipients"`
}

func (h *Handler) RetrieveForNotifications(w http.ResponseWriter, r *http.Request) {
	// Decode the JSON body into the NotificationRequest struct
	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract mentioned students from the notification text
	mentioned := extractMentionedStudents(req.Notification)

	// Get the list of students who are not suspended and are either registered to the teacher or mentioned in the notification
	var students []db.Student
	h.DB.Model(&db.Student{}).
		Select("students.email").
		Joins("LEFT JOIN teacher_students ON teacher_students.student_id = students.id").
		Joins("LEFT JOIN teachers ON teachers.id = teacher_students.teacher_id").
		Where("students.is_suspended = ?", false).
		Where("teachers.email = ? OR students.email IN ?", req.Teacher, mentioned).
		Group("students.email").
		Find(&students)

	// Prepare the list of student emails
	studentEmails := make([]string, len(students))
	for i, student := range students {
		studentEmails[i] = student.Email
	}

	// Create the response object
	response := NotificationResponse{
		Recipients: studentEmails,
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// extractMentionedStudents finds all mentioned emails in the notification text
func extractMentionedStudents(notification string) []string {
	emailRegex := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`)
	return emailRegex.FindAllString(notification, -1)
}

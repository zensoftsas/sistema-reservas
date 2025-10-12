package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"version-1-0/internal/usecase/doctor"
)

// DoctorHandler handles HTTP requests related to doctor operations
type DoctorHandler struct {
	searchDoctorsUC *doctor.SearchDoctorsUseCase
}

// NewDoctorHandler creates a new instance of DoctorHandler
func NewDoctorHandler(searchDoctorsUC *doctor.SearchDoctorsUseCase) *DoctorHandler {
	return &DoctorHandler{
		searchDoctorsUC: searchDoctorsUC,
	}
}

// Search handles the HTTP request for searching doctors
// Method: GET
// Query parameter: specialty (optional)
// Response: 200 OK with array of doctors
func (h *DoctorHandler) Search(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get specialty from query parameter (optional)
	specialty := r.URL.Query().Get("specialty")

	// Execute use case
	ctx := context.Background()
	response, err := h.searchDoctorsUC.Execute(ctx, specialty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

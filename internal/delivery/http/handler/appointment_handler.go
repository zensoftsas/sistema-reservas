package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"version-1-0/internal/delivery/http/middleware"
	"version-1-0/internal/usecase/appointment"
)

// AppointmentHandler handles HTTP requests related to appointment operations
type AppointmentHandler struct {
	createAppointmentUC *appointment.CreateAppointmentUseCase
	getByPatientUC      *appointment.GetAppointmentsByPatientUseCase
	getByDoctorUC       *appointment.GetAppointmentsByDoctorUseCase
	cancelAppointmentUC *appointment.CancelAppointmentUseCase
}

// NewAppointmentHandler creates a new instance of AppointmentHandler
func NewAppointmentHandler(
	createAppointmentUC *appointment.CreateAppointmentUseCase,
	getByPatientUC *appointment.GetAppointmentsByPatientUseCase,
	getByDoctorUC *appointment.GetAppointmentsByDoctorUseCase,
	cancelAppointmentUC *appointment.CancelAppointmentUseCase,
) *AppointmentHandler {
	return &AppointmentHandler{
		createAppointmentUC: createAppointmentUC,
		getByPatientUC:      getByPatientUC,
		getByDoctorUC:       getByDoctorUC,
		cancelAppointmentUC: cancelAppointmentUC,
	}
}

// Create handles the HTTP request for creating a new appointment
// Method: POST
// Requires: JWT token (patient creates their own appointment)
// Request body: JSON with appointment creation data
// Response: 201 Created with appointment data
func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context (patient)
	patientUserID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Decode request body
	var req appointment.CreateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.createAppointmentUC.Execute(ctx, patientUserID, req)
	if err != nil {
		// Handle specific error cases
		if err.Error() == "doctor not found" || err.Error() == "patient not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "doctor is not available at this time" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetMyAppointments handles the HTTP request for retrieving patient's own appointments
// Method: GET
// Requires: JWT token
// Response: 200 OK with array of appointments
func (h *AppointmentHandler) GetMyAppointments(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.getByPatientUC.Execute(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetDoctorAppointments handles the HTTP request for retrieving doctor's appointments
// Method: GET
// Requires: JWT token with doctor role
// Response: 200 OK with array of appointments
func (h *AppointmentHandler) GetDoctorAppointments(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get authenticated user ID from context
	doctorUserID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.getByDoctorUC.Execute(ctx, doctorUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Cancel handles the HTTP request for canceling an appointment
// Method: PUT
// Requires: JWT token (patient, doctor, or admin)
// Query parameter: id (appointment ID)
// Request body: JSON with cancellation data
// Response: 204 No Content on success
func (h *AppointmentHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is PUT
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get appointment ID from query parameter
	appointmentID := r.URL.Query().Get("id")
	if appointmentID == "" {
		http.Error(w, "Appointment ID is required", http.StatusBadRequest)
		return
	}

	// Get authenticated user info from context
	authenticatedUserID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		http.Error(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	authenticatedUserRole, ok := r.Context().Value(middleware.RoleKey).(string)
	if !ok {
		http.Error(w, "Role not found in context", http.StatusUnauthorized)
		return
	}

	// Decode request body
	var req appointment.CancelAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute use case
	ctx := context.Background()
	err := h.cancelAppointmentUC.Execute(ctx, appointmentID, authenticatedUserID, authenticatedUserRole, req)
	if err != nil {
		if err.Error() == "appointment not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "insufficient permissions to cancel this appointment" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err.Error() == "appointment is already cancelled" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response (204 No Content)
	w.WriteHeader(http.StatusNoContent)
}

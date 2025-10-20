package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"version-1-0/internal/delivery/http/middleware"
	"version-1-0/internal/usecase/appointment"
)

// AppointmentHandler handles HTTP requests related to appointment operations
type AppointmentHandler struct {
	createAppointmentUC   *appointment.CreateAppointmentUseCase
	getByPatientUC        *appointment.GetAppointmentsByPatientUseCase
	getByDoctorUC         *appointment.GetAppointmentsByDoctorUseCase
	cancelAppointmentUC   *appointment.CancelAppointmentUseCase
	confirmAppointmentUC  *appointment.ConfirmAppointmentUseCase
	completeAppointmentUC *appointment.CompleteAppointmentUseCase
	getHistoryUC          *appointment.GetPatientHistoryUseCase
}

// NewAppointmentHandler creates a new instance of AppointmentHandler
func NewAppointmentHandler(
	createAppointmentUC *appointment.CreateAppointmentUseCase,
	getByPatientUC *appointment.GetAppointmentsByPatientUseCase,
	getByDoctorUC *appointment.GetAppointmentsByDoctorUseCase,
	cancelAppointmentUC *appointment.CancelAppointmentUseCase,
	confirmAppointmentUC *appointment.ConfirmAppointmentUseCase,
	completeAppointmentUC *appointment.CompleteAppointmentUseCase,
	getHistoryUC *appointment.GetPatientHistoryUseCase,
) *AppointmentHandler {
	return &AppointmentHandler{
		createAppointmentUC:   createAppointmentUC,
		getByPatientUC:        getByPatientUC,
		getByDoctorUC:         getByDoctorUC,
		cancelAppointmentUC:   cancelAppointmentUC,
		confirmAppointmentUC:  confirmAppointmentUC,
		completeAppointmentUC: completeAppointmentUC,
		getHistoryUC:          getHistoryUC,
	}
}

// CreateAppointment godoc
// @Summary      Crear cita médica
// @Description  Crear una nueva cita médica con validaciones de disponibilidad
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        appointment  body      appointment.CreateAppointmentRequest  true  "Datos de la cita"
// @Success      201  {object}  domain.Appointment
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/appointments [post]
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

	// Validate service_id
	if req.ServiceID == "" {
		http.Error(w, "service_id is required", http.StatusBadRequest)
		return
	}

	// Parse appointment date and time
	dateTimeStr := req.AppointmentDate + " " + req.AppointmentTime + ":00"
	scheduledAt, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
	if err != nil {
		http.Error(w, "Invalid date or time format", http.StatusBadRequest)
		return
	}

	// Create appointment with service
	ctx := context.Background()
	appointmentCreated, err := h.createAppointmentUC.Execute(
		ctx,
		patientUserID,
		req.DoctorID,
		req.ServiceID,
		scheduledAt,
		req.Reason,
	)
	if err != nil {
		// Handle specific error cases
		if err.Error() == "doctor not found" || err.Error() == "patient not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "time slot is not available" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(appointmentCreated)
}

// GetMyAppointments godoc
// @Summary      Obtener mis citas
// @Description  Retorna todas las citas del usuario autenticado
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.Appointment
// @Failure      401  {object}  map[string]string
// @Router       /api/appointments/my [get]
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

// Confirm handles the HTTP request for confirming a pending appointment
// Method: PUT
// Requires: JWT token (doctor or admin)
// Query parameter: id (appointment ID)
// Response: 200 OK with confirmed appointment data
func (h *AppointmentHandler) Confirm(w http.ResponseWriter, r *http.Request) {
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

	// Execute use case
	ctx := context.Background()
	response, err := h.confirmAppointmentUC.Execute(ctx, appointmentID, authenticatedUserID, authenticatedUserRole)
	if err != nil {
		if err.Error() == "appointment not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "insufficient permissions to confirm this appointment" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err.Error() == "appointment is not in pending status" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Complete handles the HTTP request for completing a confirmed appointment
// Method: PUT
// Requires: JWT token (doctor or admin)
// Query parameter: id (appointment ID)
// Request body: JSON with completion notes
// Response: 200 OK with completed appointment data
func (h *AppointmentHandler) Complete(w http.ResponseWriter, r *http.Request) {
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
	var req appointment.CompleteAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.completeAppointmentUC.Execute(ctx, appointmentID, authenticatedUserID, authenticatedUserRole, req)
	if err != nil {
		if err.Error() == "appointment not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "insufficient permissions to complete this appointment" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		if err.Error() == "appointment is not in confirmed status" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetHistory handles the HTTP request for getting patient's medical history
// Method: GET
// Requires: JWT token
// Query parameter: patient_id (required)
// Doctors/admins can see any patient's history, patients only their own
// Response: 200 OK with array of completed appointments
func (h *AppointmentHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get patient ID from query parameter
	patientID := r.URL.Query().Get("patient_id")
	if patientID == "" {
		http.Error(w, "Patient ID is required", http.StatusBadRequest)
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

	// Execute use case
	ctx := context.Background()
	response, err := h.getHistoryUC.Execute(ctx, patientID, authenticatedUserID, authenticatedUserRole)
	if err != nil {
		if err.Error() == "patient not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "patients can only view their own medical history" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

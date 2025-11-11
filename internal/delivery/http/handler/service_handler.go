package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"version-1-0/internal/usecase/service"
)

// ServiceHandler handles HTTP requests related to service operations
type ServiceHandler struct {
	createServiceUC         *service.CreateServiceUseCase
	listServicesUC          *service.ListServicesUseCase
	assignServiceToDoctorUC *service.AssignServiceToDoctorUseCase
	getDoctorsByServiceUC   *service.GetDoctorsByServiceUseCase
	getAvailableSlotsUC     *service.GetAvailableSlotsUseCase
}

// NewServiceHandler creates a new instance of ServiceHandler
func NewServiceHandler(
	createServiceUC *service.CreateServiceUseCase,
	listServicesUC *service.ListServicesUseCase,
	assignServiceToDoctorUC *service.AssignServiceToDoctorUseCase,
	getDoctorsByServiceUC *service.GetDoctorsByServiceUseCase,
	getAvailableSlotsUC *service.GetAvailableSlotsUseCase,
) *ServiceHandler {
	return &ServiceHandler{
		createServiceUC:         createServiceUC,
		listServicesUC:          listServicesUC,
		assignServiceToDoctorUC: assignServiceToDoctorUC,
		getDoctorsByServiceUC:   getDoctorsByServiceUC,
		getAvailableSlotsUC:     getAvailableSlotsUC,
	}
}

// CreateService godoc
// @Summary      Crear servicio
// @Description  Crear un nuevo servicio médico (solo admin)
// @Tags         Services
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        service  body      dto.CreateServiceRequest  true  "Datos del servicio"
// @Success      201  {object}  dto.ServiceResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/services/create [post]
func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var req service.CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute use case
	ctx := context.Background()
	response, err := h.createServiceUC.Execute(ctx, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListServices godoc
// @Summary      Listar servicios activos
// @Description  Retorna todos los servicios médicos disponibles
// @Tags         Services
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.ServiceResponse
// @Router       /api/services [get]
func (h *ServiceHandler) List(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Execute use case
	ctx := context.Background()
	services, err := h.listServicesUC.Execute(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}

// AssignToDoctor handles the HTTP request for assigning a service to a doctor
// Method: POST
// Requires: Admin role
// Request body: JSON with doctor_id and service_id
// Response: 200 OK on success
func (h *ServiceHandler) AssignToDoctor(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var req service.AssignServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Execute use case
	ctx := context.Background()
	err := h.assignServiceToDoctorUC.Execute(ctx, req.DoctorID, req.ServiceID)
	if err != nil {
		// Handle specific error cases
		if err.Error() == "doctor not found" || err.Error() == "service not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "service is already assigned to this doctor" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Service assigned to doctor successfully",
	})
}

// GetDoctorsByService handles the HTTP request for retrieving doctors that offer a specific service
// Method: GET
// Requires: No authentication (public endpoint)
// Query parameter: service_id
// Response: 200 OK with array of doctors
func (h *ServiceHandler) GetDoctorsByService(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get service ID from query parameter
	serviceID := r.URL.Query().Get("service_id")
	if serviceID == "" {
		http.Error(w, "service_id query parameter is required", http.StatusBadRequest)
		return
	}

	// Execute use case
	ctx := context.Background()
	doctors, err := h.getDoctorsByServiceUC.Execute(ctx, serviceID)
	if err != nil {
		if err.Error() == "service not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(doctors)
}

// GetAvailableSlots godoc
// @Summary      Obtener horarios disponibles
// @Description  Calcula y retorna los slots de tiempo disponibles para un doctor, servicio y fecha específicos
// @Tags         Services
// @Accept       json
// @Produce      json
// @Param        doctor_id   query     string  true  "ID del usuario doctor"
// @Param        service_id  query     string  true  "ID del servicio"
// @Param        date        query     string  true  "Fecha en formato YYYY-MM-DD"
// @Success      200  {array}   dto.TimeSlot
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /api/services/available-slots [get]
func (h *ServiceHandler) GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	// Verify HTTP method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get query parameters
	doctorID := r.URL.Query().Get("doctor_id")
	serviceID := r.URL.Query().Get("service_id")
	dateStr := r.URL.Query().Get("date")

	// Validate parameters
	if doctorID == "" {
		http.Error(w, "doctor_id is required", http.StatusBadRequest)
		return
	}
	if serviceID == "" {
		http.Error(w, "service_id is required", http.StatusBadRequest)
		return
	}
	if dateStr == "" {
		http.Error(w, "date is required (format: YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	// Get available slots
	ctx := context.Background()
	slots, err := h.getAvailableSlotsUC.Execute(ctx, doctorID, serviceID, date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(slots)
}

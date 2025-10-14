package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"version-1-0/internal/usecase/service"
)

// ServiceHandler handles HTTP requests related to service operations
type ServiceHandler struct {
	createServiceUC         *service.CreateServiceUseCase
	listServicesUC          *service.ListServicesUseCase
	assignServiceToDoctorUC *service.AssignServiceToDoctorUseCase
	getDoctorsByServiceUC   *service.GetDoctorsByServiceUseCase
}

// NewServiceHandler creates a new instance of ServiceHandler
func NewServiceHandler(
	createServiceUC *service.CreateServiceUseCase,
	listServicesUC *service.ListServicesUseCase,
	assignServiceToDoctorUC *service.AssignServiceToDoctorUseCase,
	getDoctorsByServiceUC *service.GetDoctorsByServiceUseCase,
) *ServiceHandler {
	return &ServiceHandler{
		createServiceUC:         createServiceUC,
		listServicesUC:          listServicesUC,
		assignServiceToDoctorUC: assignServiceToDoctorUC,
		getDoctorsByServiceUC:   getDoctorsByServiceUC,
	}
}

// Create handles the HTTP request for creating a new service
// Method: POST
// Requires: Admin role
// Request body: JSON with service creation data
// Response: 201 Created with service data
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

// List handles the HTTP request for retrieving all active services
// Method: GET
// Requires: No authentication (public endpoint)
// Response: 200 OK with array of services
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

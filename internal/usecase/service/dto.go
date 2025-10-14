package service

import "time"

// CreateServiceRequest represents the input data for creating a new service
type CreateServiceRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	DurationMinutes int     `json:"duration_minutes"`
	Price           float64 `json:"price"`
}

// CreateServiceResponse represents the output data after successfully creating a service
type CreateServiceResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DurationMinutes int       `json:"duration_minutes"`
	Price           float64   `json:"price"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

// UpdateServiceRequest represents the input data for updating a service
type UpdateServiceRequest struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	DurationMinutes int     `json:"duration_minutes"`
	Price           float64 `json:"price"`
	IsActive        bool    `json:"is_active"`
}

// ServiceResponse represents a service in responses
type ServiceResponse struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	DurationMinutes int       `json:"duration_minutes"`
	Price           float64   `json:"price"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// AssignServiceRequest represents the input for assigning a service to a doctor
type AssignServiceRequest struct {
	DoctorID  string `json:"doctor_id"`
	ServiceID string `json:"service_id"`
}

// DoctorWithServicesResponse represents a doctor with their offered services
type DoctorWithServicesResponse struct {
	DoctorID   string            `json:"doctor_id"`
	DoctorName string            `json:"doctor_name"`
	Email      string            `json:"email"`
	Services   []ServiceResponse `json:"services"`
}

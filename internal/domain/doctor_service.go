package domain

import (
	"errors"
	"strings"
	"time"
)

// DoctorService represents the many-to-many relationship between doctors and services
// A doctor can offer multiple services, and a service can be offered by multiple doctors
type DoctorService struct {
	ID        string    `json:"id"`
	DoctorID  string    `json:"doctor_id"`  // References doctors.id (which references users.id)
	ServiceID string    `json:"service_id"` // References services.id
	IsActive  bool      `json:"is_active"`  // Whether the doctor currently offers this service
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate validates the doctor-service relationship data
func (ds *DoctorService) Validate() error {
	if strings.TrimSpace(ds.ID) == "" {
		return errors.New("doctor_service ID is required")
	}

	if strings.TrimSpace(ds.DoctorID) == "" {
		return errors.New("doctor ID is required")
	}

	if strings.TrimSpace(ds.ServiceID) == "" {
		return errors.New("service ID is required")
	}

	if ds.CreatedAt.IsZero() {
		return errors.New("created at is required")
	}

	if ds.UpdatedAt.IsZero() {
		return errors.New("updated at is required")
	}

	return nil
}

// Activate activates the doctor-service relationship
func (ds *DoctorService) Activate() {
	ds.IsActive = true
	ds.UpdatedAt = time.Now()
}

// Deactivate deactivates the doctor-service relationship
func (ds *DoctorService) Deactivate() {
	ds.IsActive = false
	ds.UpdatedAt = time.Now()
}

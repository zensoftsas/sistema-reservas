package domain

import (
	"errors"
	"strings"
	"time"
)

// Service represents a medical service or consultation type offered by the clinic
// Each service defines the duration (slot time) for appointments
type Service struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`                      // e.g., "Consulta Cardiológica", "Electrocardiograma"
	Description     string    `json:"description"`               // Detailed description of the service
	DurationMinutes int       `json:"duration_minutes"`          // Duration of each appointment slot (e.g., 30, 45, 60 minutes)
	Price           float64   `json:"price"`                     // Price of the service
	IsActive        bool      `json:"is_active"`                 // Whether the service is currently offered
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Validate validates the service data
func (s *Service) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return errors.New("service ID is required")
	}

	if strings.TrimSpace(s.Name) == "" {
		return errors.New("service name is required")
	}

	if s.DurationMinutes <= 0 {
		return errors.New("service duration must be greater than 0 minutes")
	}

	if s.DurationMinutes > 480 { // Max 8 hours
		return errors.New("service duration cannot exceed 480 minutes (8 hours)")
	}

	if s.Price < 0 {
		return errors.New("service price cannot be negative")
	}

	if s.CreatedAt.IsZero() {
		return errors.New("service created at is required")
	}

	if s.UpdatedAt.IsZero() {
		return errors.New("service updated at is required")
	}

	return nil
}

// Activate activates the service
func (s *Service) Activate() {
	s.IsActive = true
	s.UpdatedAt = time.Now()
}

// Deactivate deactivates the service
func (s *Service) Deactivate() {
	s.IsActive = false
	s.UpdatedAt = time.Now()
}

// CalculateEndTime calculates the end time of a service appointment
func (s *Service) CalculateEndTime(startTime time.Time) time.Time {
	return startTime.Add(time.Duration(s.DurationMinutes) * time.Minute)
}

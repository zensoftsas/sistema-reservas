package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Doctor represents a doctor entity in the medical reservation system
type Doctor struct {
	ID                 string    `json:"id"`
	UserID             string    `json:"user_id"`
	Specialty          string    `json:"specialty"`
	LicenseNumber      string    `json:"license_number"`
	YearsOfExperience  int       `json:"years_of_experience"`
	Education          string    `json:"education,omitempty"`
	Bio                string    `json:"bio,omitempty"`
	ConsultationFee    float64   `json:"consultation_fee"`
	IsAvailable        bool      `json:"is_available"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Validate checks if the Doctor entity has all required fields properly set
// Returns an error if any validation rule fails
func (d *Doctor) Validate() error {
	if strings.TrimSpace(d.ID) == "" {
		return errors.New("doctor ID is required")
	}

	if strings.TrimSpace(d.UserID) == "" {
		return errors.New("doctor user ID is required")
	}

	if strings.TrimSpace(d.Specialty) == "" {
		return errors.New("doctor specialty is required")
	}

	if strings.TrimSpace(d.LicenseNumber) == "" {
		return errors.New("doctor license number is required")
	}

	if d.YearsOfExperience < 0 {
		return errors.New("doctor years of experience cannot be negative")
	}

	if d.ConsultationFee < 0 {
		return errors.New("doctor consultation fee cannot be negative")
	}

	if d.CreatedAt.IsZero() {
		return errors.New("doctor created at is required")
	}

	if d.UpdatedAt.IsZero() {
		return errors.New("doctor updated at is required")
	}

	return nil
}

// FormattedFee returns the consultation fee formatted as a string with currency symbol
// Uses the format "S/ X.XX" (Peruvian Sol) for the clinic
func (d *Doctor) FormattedFee() string {
	return fmt.Sprintf("S/ %.2f", d.ConsultationFee)
}

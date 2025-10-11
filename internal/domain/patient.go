package domain

import (
	"errors"
	"strings"
	"time"
)

// Patient represents a patient entity in the medical reservation system
type Patient struct {
	ID                    string    `json:"id"`
	UserID                string    `json:"user_id"`
	Birthdate             time.Time `json:"birthdate"`
	DocumentType          string    `json:"document_type"`
	DocumentNumber        string    `json:"document_number"`
	Address               string    `json:"address"`
	EmergencyContactName  string    `json:"emergency_contact_name"`
	EmergencyContactPhone string    `json:"emergency_contact_phone"`
	BloodType             string    `json:"blood_type,omitempty"`
	Allergies             []string  `json:"allergies,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// Validate checks if the Patient entity has all required fields properly set
// Returns an error if any validation rule fails
func (p *Patient) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return errors.New("patient ID is required")
	}

	if strings.TrimSpace(p.UserID) == "" {
		return errors.New("patient user ID is required")
	}

	if p.Birthdate.IsZero() {
		return errors.New("patient birthdate is required")
	}

	if p.Birthdate.After(time.Now()) {
		return errors.New("patient birthdate cannot be in the future")
	}

	if strings.TrimSpace(p.DocumentType) == "" {
		return errors.New("patient document type is required")
	}

	if strings.TrimSpace(p.DocumentNumber) == "" {
		return errors.New("patient document number is required")
	}

	if strings.TrimSpace(p.Address) == "" {
		return errors.New("patient address is required")
	}

	if strings.TrimSpace(p.EmergencyContactName) == "" {
		return errors.New("patient emergency contact name is required")
	}

	if strings.TrimSpace(p.EmergencyContactPhone) == "" {
		return errors.New("patient emergency contact phone is required")
	}

	if p.CreatedAt.IsZero() {
		return errors.New("patient created at is required")
	}

	if p.UpdatedAt.IsZero() {
		return errors.New("patient updated at is required")
	}

	return nil
}

// Age calculates and returns the patient's current age in years based on their birthdate
func (p *Patient) Age() int {
	now := time.Now()
	age := now.Year() - p.Birthdate.Year()

	// Adjust age if birthday hasn't occurred yet this year
	if now.Month() < p.Birthdate.Month() ||
		(now.Month() == p.Birthdate.Month() && now.Day() < p.Birthdate.Day()) {
		age--
	}

	return age
}

package domain

import (
	"errors"
	"strings"
	"time"
)

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

// Appointment status constants
const (
	StatusPending   AppointmentStatus = "pending"
	StatusConfirmed AppointmentStatus = "confirmed"
	StatusCancelled AppointmentStatus = "cancelled"
	StatusCompleted AppointmentStatus = "completed"
)

// Appointment represents an appointment entity in the medical reservation system
type Appointment struct {
	ID                 string            `json:"id"`
	PatientID          string            `json:"patient_id"`
	DoctorID           string            `json:"doctor_id"`
	ScheduledAt        time.Time         `json:"scheduled_at"`
	Duration           int               `json:"duration"` // in minutes
	Reason             string            `json:"reason"`
	Notes              string            `json:"notes"`
	Status             AppointmentStatus `json:"status"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
	CancelledAt        *time.Time        `json:"cancelled_at,omitempty"`
	CancellationReason string            `json:"cancellation_reason,omitempty"`
}

// Validate checks if the Appointment entity has all required fields properly set
// Returns an error if any validation rule fails
func (a *Appointment) Validate() error {
	if strings.TrimSpace(a.ID) == "" {
		return errors.New("appointment ID is required")
	}

	if strings.TrimSpace(a.PatientID) == "" {
		return errors.New("appointment patient ID is required")
	}

	if strings.TrimSpace(a.DoctorID) == "" {
		return errors.New("appointment doctor ID is required")
	}

	if a.ScheduledAt.IsZero() {
		return errors.New("appointment scheduled time is required")
	}

	if a.ScheduledAt.Before(time.Now()) {
		return errors.New("appointment cannot be scheduled in the past")
	}

	if a.Duration <= 0 {
		return errors.New("appointment duration must be greater than 0")
	}

	if strings.TrimSpace(a.Reason) == "" {
		return errors.New("appointment reason is required")
	}

	if a.Status != StatusPending && a.Status != StatusConfirmed &&
		a.Status != StatusCancelled && a.Status != StatusCompleted {
		return errors.New("invalid appointment status")
	}

	if a.CreatedAt.IsZero() {
		return errors.New("appointment created at is required")
	}

	if a.UpdatedAt.IsZero() {
		return errors.New("appointment updated at is required")
	}

	return nil
}

// CanBeCancelled verifies if the appointment can be cancelled
// Returns an error if the appointment cannot be cancelled
func (a *Appointment) CanBeCancelled() error {
	if a.Status == StatusCancelled {
		return errors.New("appointment is already cancelled")
	}

	if a.Status == StatusCompleted {
		return errors.New("completed appointment cannot be cancelled")
	}

	hoursUntilAppointment := time.Until(a.ScheduledAt).Hours()
	if hoursUntilAppointment < 24 {
		return errors.New("appointment must be cancelled at least 24 hours in advance")
	}

	return nil
}

// Cancel cancels the appointment with the given reason
// Returns an error if the appointment cannot be cancelled
func (a *Appointment) Cancel(reason string) error {
	if err := a.CanBeCancelled(); err != nil {
		return err
	}

	if strings.TrimSpace(reason) == "" {
		return errors.New("cancellation reason is required")
	}

	now := time.Now()
	a.Status = StatusCancelled
	a.CancelledAt = &now
	a.CancellationReason = reason
	a.UpdatedAt = now

	return nil
}

// Confirm confirms the appointment
// Returns an error if the appointment is not in pending status
func (a *Appointment) Confirm() error {
	if a.Status != StatusPending {
		return errors.New("only pending appointments can be confirmed")
	}

	if a.ScheduledAt.Before(time.Now()) {
		return errors.New("cannot confirm an appointment scheduled in the past")
	}

	a.Status = StatusConfirmed
	a.UpdatedAt = time.Now()

	return nil
}

// Complete marks the appointment as completed with optional notes
// Returns an error if the appointment is not in confirmed status
func (a *Appointment) Complete(notes string) error {
	if a.Status != StatusConfirmed {
		return errors.New("only confirmed appointments can be completed")
	}

	a.Status = StatusCompleted
	a.Notes = notes
	a.UpdatedAt = time.Now()

	return nil
}

// IsPast returns true if the appointment's scheduled time has passed
func (a *Appointment) IsPast() bool {
	return a.ScheduledAt.Before(time.Now())
}

// EndTime calculates and returns the appointment's end time
// based on the scheduled time and duration
func (a *Appointment) EndTime() time.Time {
	return a.ScheduledAt.Add(time.Duration(a.Duration) * time.Minute)
}

// IsValidStatus checks if a given status string is a valid AppointmentStatus
func IsValidAppointmentStatus(status string) bool {
	s := AppointmentStatus(status)
	return s == StatusPending || s == StatusConfirmed ||
		s == StatusCancelled || s == StatusCompleted
}

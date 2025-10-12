package appointment

import (
	"context"
	"errors"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// CancelAppointmentUseCase handles the business logic for canceling appointments
type CancelAppointmentUseCase struct {
	appointmentRepo repository.AppointmentRepository
}

// NewCancelAppointmentUseCase creates a new instance of CancelAppointmentUseCase
func NewCancelAppointmentUseCase(appointmentRepo repository.AppointmentRepository) *CancelAppointmentUseCase {
	return &CancelAppointmentUseCase{
		appointmentRepo: appointmentRepo,
	}
}

// Execute cancels an appointment with permission validation
// Only the patient, the doctor involved, or an admin can cancel an appointment
func (uc *CancelAppointmentUseCase) Execute(ctx context.Context, appointmentID string, authenticatedUserID string, authenticatedUserRole string, req CancelAppointmentRequest) error {
	// Retrieve the appointment
	appointment, err := uc.appointmentRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return err
	}
	if appointment == nil {
		return errors.New("appointment not found")
	}

	// Verify permissions: only the patient, the doctor, or an admin can cancel
	if authenticatedUserRole != "admin" &&
		authenticatedUserID != appointment.PatientID &&
		authenticatedUserID != appointment.DoctorID {
		return errors.New("insufficient permissions to cancel this appointment")
	}

	// Verify that the appointment is not already cancelled
	if appointment.Status == domain.StatusCancelled {
		return errors.New("appointment is already cancelled")
	}

	// Update appointment status and notes
	appointment.Status = domain.StatusCancelled

	// Build cancellation note
	cancellationNote := "Cancelled"
	if req.Reason != "" {
		cancellationNote += ": " + req.Reason
	}
	appointment.Notes = cancellationNote
	appointment.UpdatedAt = time.Now()

	// Save changes to database
	err = uc.appointmentRepo.Update(ctx, appointment)
	if err != nil {
		return err
	}

	return nil
}

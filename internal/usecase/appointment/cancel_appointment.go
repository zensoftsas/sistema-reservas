package appointment

import (
	"context"
	"errors"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
	"version-1-0/pkg/email"
)

// CancelAppointmentUseCase handles the business logic for canceling appointments
type CancelAppointmentUseCase struct {
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository
	emailService    *email.EmailService
}

// NewCancelAppointmentUseCase creates a new instance of CancelAppointmentUseCase
func NewCancelAppointmentUseCase(appointmentRepo repository.AppointmentRepository, userRepo repository.UserRepository, emailService *email.EmailService) *CancelAppointmentUseCase {
	return &CancelAppointmentUseCase{
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
		emailService:    emailService,
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

	// Get patient and doctor info for email
	patient, _ := uc.userRepo.FindByID(ctx, appointment.PatientID)
	doctor, _ := uc.userRepo.FindByID(ctx, appointment.DoctorID)

	// Send email notifications
	if uc.emailService != nil && patient != nil && doctor != nil {
		patientName := patient.FirstName + " " + patient.LastName
		doctorName := doctor.FirstName + " " + doctor.LastName
		date := appointment.ScheduledAt.Format("2006-01-02")
		time := appointment.ScheduledAt.Format("15:04")

		// Notify patient
		uc.emailService.SendAppointmentCancelled(
			patient.Email,
			patientName,
			doctorName,
			date,
			time,
			req.Reason,
		)

		// Notify doctor
		uc.emailService.SendAppointmentCancelled(
			doctor.Email,
			doctorName,
			doctorName,
			date,
			time,
			req.Reason,
		)
	}

	return nil
}

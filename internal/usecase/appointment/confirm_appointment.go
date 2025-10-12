package appointment

import (
	"context"
	"errors"

	"version-1-0/internal/repository"
)

// ConfirmAppointmentUseCase handles the business logic for confirming appointments
type ConfirmAppointmentUseCase struct {
	appointmentRepo repository.AppointmentRepository
}

// NewConfirmAppointmentUseCase creates a new instance of ConfirmAppointmentUseCase
func NewConfirmAppointmentUseCase(appointmentRepo repository.AppointmentRepository) *ConfirmAppointmentUseCase {
	return &ConfirmAppointmentUseCase{
		appointmentRepo: appointmentRepo,
	}
}

// Execute confirms an appointment
// Only doctors or admins can confirm appointments
func (uc *ConfirmAppointmentUseCase) Execute(ctx context.Context, appointmentID string, authenticatedUserID string, authenticatedUserRole string) (*ConfirmAppointmentResponse, error) {
	// Retrieve the appointment
	appointment, err := uc.appointmentRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	if appointment == nil {
		return nil, errors.New("appointment not found")
	}

	// Verify permissions: only the doctor or an admin can confirm
	if authenticatedUserRole != "admin" && authenticatedUserID != appointment.DoctorID {
		return nil, errors.New("insufficient permissions to confirm this appointment")
	}

	// Use domain method to confirm
	err = appointment.Confirm()
	if err != nil {
		return nil, err
	}

	// Save changes to database
	err = uc.appointmentRepo.Update(ctx, appointment)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &ConfirmAppointmentResponse{
		ID:              appointment.ID,
		PatientID:       appointment.PatientID,
		DoctorID:        appointment.DoctorID,
		AppointmentDate: appointment.ScheduledAt.Format("2006-01-02"),
		AppointmentTime: appointment.ScheduledAt.Format("15:04"),
		Status:          string(appointment.Status),
		Reason:          appointment.Reason,
		UpdatedAt:       appointment.UpdatedAt,
	}

	return response, nil
}

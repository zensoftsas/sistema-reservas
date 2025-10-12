package appointment

import (
	"context"
	"errors"

	"version-1-0/internal/repository"
)

// CompleteAppointmentUseCase handles the business logic for completing appointments
type CompleteAppointmentUseCase struct {
	appointmentRepo repository.AppointmentRepository
}

// NewCompleteAppointmentUseCase creates a new instance of CompleteAppointmentUseCase
func NewCompleteAppointmentUseCase(appointmentRepo repository.AppointmentRepository) *CompleteAppointmentUseCase {
	return &CompleteAppointmentUseCase{
		appointmentRepo: appointmentRepo,
	}
}

// Execute completes an appointment with medical notes
// Only doctors or admins can complete appointments
func (uc *CompleteAppointmentUseCase) Execute(ctx context.Context, appointmentID string, authenticatedUserID string, authenticatedUserRole string, req CompleteAppointmentRequest) (*CompleteAppointmentResponse, error) {
	// Retrieve the appointment
	appointment, err := uc.appointmentRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return nil, err
	}
	if appointment == nil {
		return nil, errors.New("appointment not found")
	}

	// Verify permissions: only the doctor or an admin can complete
	if authenticatedUserRole != "admin" && authenticatedUserID != appointment.DoctorID {
		return nil, errors.New("insufficient permissions to complete this appointment")
	}

	// Use domain method to complete with notes
	err = appointment.Complete(req.Notes)
	if err != nil {
		return nil, err
	}

	// Save changes to database
	err = uc.appointmentRepo.Update(ctx, appointment)
	if err != nil {
		return nil, err
	}

	// Build and return response
	response := &CompleteAppointmentResponse{
		ID:              appointment.ID,
		PatientID:       appointment.PatientID,
		DoctorID:        appointment.DoctorID,
		AppointmentDate: appointment.ScheduledAt.Format("2006-01-02"),
		AppointmentTime: appointment.ScheduledAt.Format("15:04"),
		Status:          string(appointment.Status),
		Reason:          appointment.Reason,
		Notes:           appointment.Notes,
		UpdatedAt:       appointment.UpdatedAt,
	}

	return response, nil
}

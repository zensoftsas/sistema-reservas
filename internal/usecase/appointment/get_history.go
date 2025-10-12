package appointment

import (
	"context"
	"errors"

	"version-1-0/internal/repository"
)

// GetPatientHistoryUseCase handles retrieving completed appointments (medical history)
type GetPatientHistoryUseCase struct {
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository
}

// NewGetPatientHistoryUseCase creates a new instance
func NewGetPatientHistoryUseCase(appointmentRepo repository.AppointmentRepository, userRepo repository.UserRepository) *GetPatientHistoryUseCase {
	return &GetPatientHistoryUseCase{
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
	}
}

// Execute retrieves medical history (completed appointments) for a patient
// Doctors and admins can see any patient's history
// Patients can only see their own history
func (uc *GetPatientHistoryUseCase) Execute(ctx context.Context, patientID string, authenticatedUserID string, authenticatedUserRole string) ([]GetAppointmentResponse, error) {
	// Verify permissions
	if authenticatedUserRole == "patient" && authenticatedUserID != patientID {
		return nil, errors.New("patients can only view their own medical history")
	}

	// Verify patient exists
	patient, err := uc.userRepo.FindByID(ctx, patientID)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, errors.New("patient not found")
	}

	// Get all appointments for patient
	appointments, err := uc.appointmentRepo.FindByPatientID(ctx, patientID)
	if err != nil {
		return nil, err
	}

	// Filter only completed appointments and convert to response
	var history []GetAppointmentResponse
	for _, appointment := range appointments {
		// Only include completed appointments in history
		if appointment.Status == "completed" {
			history = append(history, GetAppointmentResponse{
				ID:              appointment.ID,
				PatientID:       appointment.PatientID,
				DoctorID:        appointment.DoctorID,
				AppointmentDate: appointment.ScheduledAt.Format("2006-01-02"),
				AppointmentTime: appointment.ScheduledAt.Format("15:04"),
				Status:          string(appointment.Status),
				Reason:          appointment.Reason,
				Notes:           appointment.Notes,
				CreatedAt:       appointment.CreatedAt,
			})
		}
	}

	return history, nil
}

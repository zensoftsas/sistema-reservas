package appointment

import (
	"context"

	"version-1-0/internal/repository"
)

// GetAppointmentsByPatientUseCase handles retrieving all appointments for a patient
type GetAppointmentsByPatientUseCase struct {
	appointmentRepo repository.AppointmentRepository
}

// NewGetAppointmentsByPatientUseCase creates a new instance of GetAppointmentsByPatientUseCase
func NewGetAppointmentsByPatientUseCase(appointmentRepo repository.AppointmentRepository) *GetAppointmentsByPatientUseCase {
	return &GetAppointmentsByPatientUseCase{
		appointmentRepo: appointmentRepo,
	}
}

// Execute retrieves all appointments for a specific patient
func (uc *GetAppointmentsByPatientUseCase) Execute(ctx context.Context, patientUserID string) ([]GetAppointmentResponse, error) {
	// Retrieve appointments from repository
	appointments, err := uc.appointmentRepo.FindByPatientID(ctx, patientUserID)
	if err != nil {
		return nil, err
	}

	// Convert domain appointments to response DTOs
	responses := make([]GetAppointmentResponse, len(appointments))
	for i, appointment := range appointments {
		responses[i] = GetAppointmentResponse{
			ID:              appointment.ID,
			PatientID:       appointment.PatientID,
			DoctorID:        appointment.DoctorID,
			AppointmentDate: appointment.ScheduledAt.Format("2006-01-02"),
			AppointmentTime: appointment.ScheduledAt.Format("15:04"),
			Status:          string(appointment.Status),
			Reason:          appointment.Reason,
			Notes:           appointment.Notes,
			CreatedAt:       appointment.CreatedAt,
		}
	}

	return responses, nil
}

// GetAppointmentsByDoctorUseCase handles retrieving all appointments for a doctor
type GetAppointmentsByDoctorUseCase struct {
	appointmentRepo repository.AppointmentRepository
}

// NewGetAppointmentsByDoctorUseCase creates a new instance of GetAppointmentsByDoctorUseCase
func NewGetAppointmentsByDoctorUseCase(appointmentRepo repository.AppointmentRepository) *GetAppointmentsByDoctorUseCase {
	return &GetAppointmentsByDoctorUseCase{
		appointmentRepo: appointmentRepo,
	}
}

// Execute retrieves all appointments for a specific doctor
func (uc *GetAppointmentsByDoctorUseCase) Execute(ctx context.Context, doctorUserID string) ([]GetAppointmentResponse, error) {
	// Retrieve appointments from repository
	appointments, err := uc.appointmentRepo.FindByDoctorID(ctx, doctorUserID)
	if err != nil {
		return nil, err
	}

	// Convert domain appointments to response DTOs
	responses := make([]GetAppointmentResponse, len(appointments))
	for i, appointment := range appointments {
		responses[i] = GetAppointmentResponse{
			ID:              appointment.ID,
			PatientID:       appointment.PatientID,
			DoctorID:        appointment.DoctorID,
			AppointmentDate: appointment.ScheduledAt.Format("2006-01-02"),
			AppointmentTime: appointment.ScheduledAt.Format("15:04"),
			Status:          string(appointment.Status),
			Reason:          appointment.Reason,
			Notes:           appointment.Notes,
			CreatedAt:       appointment.CreatedAt,
		}
	}

	return responses, nil
}

package appointment

import (
	"context"
	"errors"

	"version-1-0/internal/repository"
)

// GetAppointmentsByPatientUseCase handles retrieving all appointments for a patient
type GetAppointmentsByPatientUseCase struct {
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository
}

// NewGetAppointmentsByPatientUseCase creates a new instance of GetAppointmentsByPatientUseCase
func NewGetAppointmentsByPatientUseCase(appointmentRepo repository.AppointmentRepository, userRepo repository.UserRepository) *GetAppointmentsByPatientUseCase {
	return &GetAppointmentsByPatientUseCase{
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
	}
}

// Execute retrieves all appointments for a specific patient
func (uc *GetAppointmentsByPatientUseCase) Execute(ctx context.Context, patientUserID string) ([]GetAppointmentResponse, error) {
	// Convert user_id to patient.id from patients table
	patientID, err := uc.userRepo.FindPatientIDByUserID(ctx, patientUserID)
	if err != nil {
		return nil, errors.New("patient profile not found for user")
	}

	// Retrieve appointments from repository using the real patient.id
	appointments, err := uc.appointmentRepo.FindByPatientID(ctx, patientID)
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
			PatientName:     appointment.PatientName,
			DoctorName:      appointment.DoctorName,
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
	userRepo        repository.UserRepository
}

// NewGetAppointmentsByDoctorUseCase creates a new instance of GetAppointmentsByDoctorUseCase
func NewGetAppointmentsByDoctorUseCase(appointmentRepo repository.AppointmentRepository, userRepo repository.UserRepository) *GetAppointmentsByDoctorUseCase {
	return &GetAppointmentsByDoctorUseCase{
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
	}
}

// Execute retrieves all appointments for a specific doctor
func (uc *GetAppointmentsByDoctorUseCase) Execute(ctx context.Context, doctorUserID string) ([]GetAppointmentResponse, error) {
	// Convert user_id to doctor.id from doctors table
	doctorID, err := uc.userRepo.FindDoctorIDByUserID(ctx, doctorUserID)
	if err != nil {
		return nil, errors.New("doctor profile not found for user")
	}

	// Retrieve appointments from repository using the real doctor.id
	appointments, err := uc.appointmentRepo.FindByDoctorID(ctx, doctorID)
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
			PatientName:     appointment.PatientName,
			DoctorName:      appointment.DoctorName,
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

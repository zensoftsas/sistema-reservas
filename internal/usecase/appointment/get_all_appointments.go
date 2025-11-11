package appointment

import (
	"context"
	"time"

	"version-1-0/internal/repository"
)

// GetAllAppointmentsUseCase handles retrieving all appointments with filters (admin only)
type GetAllAppointmentsUseCase struct {
	appointmentRepo repository.AppointmentRepository
}

// NewGetAllAppointmentsUseCase creates a new instance of GetAllAppointmentsUseCase
func NewGetAllAppointmentsUseCase(appointmentRepo repository.AppointmentRepository) *GetAllAppointmentsUseCase {
	return &GetAllAppointmentsUseCase{
		appointmentRepo: appointmentRepo,
	}
}

// Execute retrieves all appointments with optional filters
func (uc *GetAllAppointmentsUseCase) Execute(ctx context.Context, req GetAllAppointmentsRequest) ([]GetAppointmentResponse, error) {
	// Build filters
	filters := repository.AppointmentFilters{
		Status:    req.Status,
		DoctorID:  req.DoctorID,
		PatientID: req.PatientID,
		ServiceID: req.ServiceID,
	}

	// Parse date filters if provided
	if req.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
		if err == nil {
			filters.DateFrom = &dateFrom
		}
	}

	if req.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", req.DateTo)
		if err == nil {
			// Set to end of day
			endOfDay := dateTo.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filters.DateTo = &endOfDay
		}
	}

	// Retrieve appointments from repository with filters
	appointments, err := uc.appointmentRepo.FindAllWithFilters(ctx, filters)
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
			ServiceID:       appointment.ServiceID,
			PatientName:     appointment.PatientName,
			DoctorName:      appointment.DoctorName,
			ServiceName:     appointment.ServiceName,
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

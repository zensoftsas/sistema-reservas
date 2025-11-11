package appointment

import (
	"context"
	"errors"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// RescheduleAppointmentUseCase handles rescheduling an appointment (admin only)
type RescheduleAppointmentUseCase struct {
	appointmentRepo repository.AppointmentRepository
	serviceRepo     repository.ServiceRepository
}

// NewRescheduleAppointmentUseCase creates a new instance of RescheduleAppointmentUseCase
func NewRescheduleAppointmentUseCase(
	appointmentRepo repository.AppointmentRepository,
	serviceRepo repository.ServiceRepository,
) *RescheduleAppointmentUseCase {
	return &RescheduleAppointmentUseCase{
		appointmentRepo: appointmentRepo,
		serviceRepo:     serviceRepo,
	}
}

// Execute reschedules an appointment to a new date/time with validations
func (uc *RescheduleAppointmentUseCase) Execute(
	ctx context.Context,
	appointmentID string,
	authenticatedUserID string,
	authenticatedUserRole string,
	req RescheduleAppointmentRequest,
) (*RescheduleAppointmentResponse, error) {
	// Find the appointment
	appointment, err := uc.appointmentRepo.FindByID(ctx, appointmentID)
	if err != nil {
		return nil, errors.New("appointment not found")
	}

	// Check permissions: patients can only reschedule their own appointments, admins can reschedule any
	if authenticatedUserRole != "admin" {
		if appointment.PatientID != authenticatedUserID {
			return nil, errors.New("insufficient permissions to reschedule this appointment")
		}
	}

	// Check if appointment can be rescheduled
	if appointment.Status == domain.StatusCancelled {
		return nil, errors.New("cannot reschedule a cancelled appointment")
	}

	if appointment.Status == domain.StatusCompleted {
		return nil, errors.New("cannot reschedule a completed appointment")
	}

	// Parse new scheduled time in Peru timezone (America/Lima UTC-5)
	location, err := time.LoadLocation("America/Lima")
	if err != nil {
		// Fallback to UTC-5 if location loading fails
		location = time.FixedZone("America/Lima", -5*60*60)
	}

	dateTimeStr := req.NewDate + " " + req.NewTime + ":00"
	newScheduledAt, err := time.ParseInLocation("2006-01-02 15:04:05", dateTimeStr, location)
	if err != nil {
		return nil, errors.New("invalid date or time format")
	}

	// Validate new time is in the future
	if newScheduledAt.Before(time.Now()) {
		return nil, errors.New("new appointment time must be in the future")
	}

	// Get service to know the duration
	var duration int
	if appointment.ServiceID != "" {
		service, err := uc.serviceRepo.FindByID(ctx, appointment.ServiceID)
		if err == nil && service != nil {
			duration = service.DurationMinutes
		} else {
			// Use appointment's existing duration if service not found
			duration = appointment.Duration
		}
	} else {
		duration = appointment.Duration
	}

	// Check for conflicts with other appointments
	endTime := newScheduledAt.Add(time.Duration(duration) * time.Minute)

	// Find all appointments for this doctor on the new date
	existingAppointments, err := uc.appointmentRepo.FindByDoctorAndDate(ctx, appointment.DoctorID, newScheduledAt)
	if err != nil {
		return nil, errors.New("failed to check doctor availability")
	}

	// Check for time slot conflicts (exclude current appointment)
	for _, existing := range existingAppointments {
		// Skip the appointment being rescheduled
		if existing.ID == appointmentID {
			continue
		}

		// Skip cancelled appointments
		if existing.Status == domain.StatusCancelled {
			continue
		}

		existingStart := existing.ScheduledAt
		existingEnd := existing.EndTime()

		// Check if time slots overlap
		if (newScheduledAt.Before(existingEnd) && endTime.After(existingStart)) {
			return nil, errors.New("time slot conflicts with another appointment")
		}
	}

	// Update appointment with new scheduled time
	appointment.ScheduledAt = newScheduledAt
	appointment.UpdatedAt = time.Now()

	// Save updated appointment
	if err := uc.appointmentRepo.Update(ctx, appointment); err != nil {
		return nil, errors.New("failed to reschedule appointment")
	}

	// Return response
	response := &RescheduleAppointmentResponse{
		ID:              appointment.ID,
		PatientID:       appointment.PatientID,
		DoctorID:        appointment.DoctorID,
		PatientName:     appointment.PatientName,
		DoctorName:      appointment.DoctorName,
		ServiceName:     appointment.ServiceName,
		AppointmentDate: appointment.ScheduledAt.Format("2006-01-02"),
		AppointmentTime: appointment.ScheduledAt.Format("15:04"),
		Status:          string(appointment.Status),
		Reason:          appointment.Reason,
		UpdatedAt:       appointment.UpdatedAt,
	}

	return response, nil
}

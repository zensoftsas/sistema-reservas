package appointment

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
	"version-1-0/pkg/email"
)

// CreateAppointmentUseCase handles the business logic for creating appointments
type CreateAppointmentUseCase struct {
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository
	emailService    *email.EmailService
}

// NewCreateAppointmentUseCase creates a new instance of CreateAppointmentUseCase
func NewCreateAppointmentUseCase(appointmentRepo repository.AppointmentRepository, userRepo repository.UserRepository, emailService *email.EmailService) *CreateAppointmentUseCase {
	return &CreateAppointmentUseCase{
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
		emailService:    emailService,
	}
}

// Execute creates a new appointment with comprehensive validation
func (uc *CreateAppointmentUseCase) Execute(ctx context.Context, patientUserID string, req CreateAppointmentRequest) (*CreateAppointmentResponse, error) {
	// Validate doctor ID is not empty
	if strings.TrimSpace(req.DoctorID) == "" {
		return nil, errors.New("doctor ID is required")
	}

	// Validate appointment date is not empty
	if strings.TrimSpace(req.AppointmentDate) == "" {
		return nil, errors.New("appointment date is required")
	}

	// Validate appointment time is not empty
	if strings.TrimSpace(req.AppointmentTime) == "" {
		return nil, errors.New("appointment time is required")
	}

	// Validate reason is not empty and has minimum length
	if strings.TrimSpace(req.Reason) == "" {
		return nil, errors.New("reason is required")
	}
	if len(strings.TrimSpace(req.Reason)) < 10 {
		return nil, errors.New("reason must be at least 10 characters long")
	}

	// Combine date and time into a single datetime
	dateTimeStr := req.AppointmentDate + " " + req.AppointmentTime + ":00"
	scheduledAt, err := time.Parse("2006-01-02 15:04:05", dateTimeStr)
	if err != nil {
		return nil, errors.New("invalid date or time format")
	}

	// Validate that the appointment is in the future
	if scheduledAt.Before(time.Now()) {
		return nil, errors.New("appointment must be scheduled in the future")
	}

	// Verify that the patient exists and is active
	patient, err := uc.userRepo.FindByID(ctx, patientUserID)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, errors.New("patient not found")
	}
	if !patient.IsActive {
		return nil, errors.New("patient account is inactive")
	}

	// Verify that the doctor exists and is active
	doctor, err := uc.userRepo.FindByID(ctx, req.DoctorID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}
	if !doctor.IsActive {
		return nil, errors.New("doctor account is inactive")
	}
	if doctor.Role != "doctor" {
		return nil, errors.New("user is not a doctor")
	}

	// Verify doctor availability at the requested time
	existingAppointments, err := uc.appointmentRepo.FindByDoctorAndDate(ctx, req.DoctorID, scheduledAt)
	if err != nil {
		return nil, err
	}

	// Check for time conflicts
	for _, existing := range existingAppointments {
		if existing.ScheduledAt.Format("15:04") == scheduledAt.Format("15:04") &&
			existing.Status != domain.StatusCancelled {
			return nil, errors.New("doctor is not available at this time")
		}
	}

	// Create the appointment
	appointment := &domain.Appointment{
		ID:          uuid.New().String(),
		PatientID:   patientUserID,
		DoctorID:    req.DoctorID,
		ScheduledAt: scheduledAt,
		Duration:    30, // Default 30 minutes
		Reason:      strings.TrimSpace(req.Reason),
		Notes:       "",
		Status:      domain.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	err = uc.appointmentRepo.Create(ctx, appointment)
	if err != nil {
		return nil, err
	}

	// Build and return response
	response := &CreateAppointmentResponse{
		ID:              appointment.ID,
		PatientID:       appointment.PatientID,
		DoctorID:        appointment.DoctorID,
		AppointmentDate: appointment.ScheduledAt.Format("2006-01-02"),
		AppointmentTime: appointment.ScheduledAt.Format("15:04"),
		Status:          string(appointment.Status),
		Reason:          appointment.Reason,
		CreatedAt:       appointment.CreatedAt,
	}

	// Send email notification to patient
	if uc.emailService != nil {
		patientName := patient.FirstName + " " + patient.LastName
		doctorName := doctor.FirstName + " " + doctor.LastName
		uc.emailService.SendAppointmentCreated(
			patient.Email,
			patientName,
			doctorName,
			response.AppointmentDate,
			response.AppointmentTime,
		)
	}

	return response, nil
}

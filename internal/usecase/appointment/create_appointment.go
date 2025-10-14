package appointment

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
	"version-1-0/pkg/email"
)

// CreateAppointmentUseCase handles the creation of appointments
type CreateAppointmentUseCase struct {
	appointmentRepo   repository.AppointmentRepository
	userRepo          repository.UserRepository
	serviceRepo       repository.ServiceRepository
	doctorServiceRepo repository.DoctorServiceRepository
	emailService      *email.EmailService
}

// NewCreateAppointmentUseCase creates a new CreateAppointmentUseCase
func NewCreateAppointmentUseCase(
	appointmentRepo repository.AppointmentRepository,
	userRepo repository.UserRepository,
	serviceRepo repository.ServiceRepository,
	doctorServiceRepo repository.DoctorServiceRepository,
	emailService *email.EmailService,
) *CreateAppointmentUseCase {
	return &CreateAppointmentUseCase{
		appointmentRepo:   appointmentRepo,
		userRepo:          userRepo,
		serviceRepo:       serviceRepo,
		doctorServiceRepo: doctorServiceRepo,
		emailService:      emailService,
	}
}

// Execute creates a new appointment with a service
func (uc *CreateAppointmentUseCase) Execute(ctx context.Context, patientID, doctorID, serviceID string, scheduledAt time.Time, reason string) (*domain.Appointment, error) {
	// Validate patient exists
	patient, err := uc.userRepo.FindByID(ctx, patientID)
	if err != nil {
		return nil, err
	}
	if patient == nil {
		return nil, errors.New("patient not found")
	}

	// Validate doctor exists
	doctor, err := uc.userRepo.FindByID(ctx, doctorID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}

	// Get real doctor.id
	realDoctorID, err := uc.userRepo.FindDoctorIDByUserID(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	// Validate service exists
	service, err := uc.serviceRepo.FindByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, errors.New("service not found")
	}
	if !service.IsActive {
		return nil, errors.New("service is not active")
	}

	// Validate doctor offers this service
	isAssigned, err := uc.doctorServiceRepo.IsAssigned(ctx, realDoctorID, serviceID)
	if err != nil {
		return nil, err
	}
	if !isAssigned {
		return nil, errors.New("doctor does not offer this service")
	}

	// Check for scheduling conflicts
	appointmentEnd := scheduledAt.Add(time.Duration(service.DurationMinutes) * time.Minute)
	conflicts, err := uc.appointmentRepo.FindByDoctorAndDate(ctx, realDoctorID, scheduledAt)
	if err != nil {
		return nil, err
	}

	for _, conflict := range conflicts {
		if conflict.Status == "cancelled" {
			continue
		}

		conflictEnd := conflict.ScheduledAt.Add(time.Duration(conflict.Duration) * time.Minute)

		if scheduledAt.Before(conflictEnd) && appointmentEnd.After(conflict.ScheduledAt) {
			return nil, errors.New("time slot is not available")
		}
	}

	// Create appointment
	now := time.Now()
	appointment := &domain.Appointment{
		ID:          uuid.New().String(),
		PatientID:   patientID,
		DoctorID:    realDoctorID,
		ServiceID:   serviceID,
		ServiceName: service.Name,
		ScheduledAt: scheduledAt,
		Duration:    service.DurationMinutes,
		Reason:      reason,
		Status:      "pending",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := appointment.Validate(); err != nil {
		return nil, err
	}

	if err := uc.appointmentRepo.Create(ctx, appointment); err != nil {
		return nil, err
	}

	// Send notification email
	if uc.emailService != nil {
		patientName := patient.FirstName + " " + patient.LastName
		doctorName := doctor.FirstName + " " + doctor.LastName
		date := scheduledAt.Format("2006-01-02")
		timeStr := scheduledAt.Format("15:04")

		go func() {
			if err := uc.emailService.SendAppointmentCreated(patient.Email, patientName, doctorName, date, timeStr); err != nil {
				log.Printf("Failed to send appointment created email: %v", err)
			}
		}()
	}

	return appointment, nil
}

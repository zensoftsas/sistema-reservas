package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// AssignServiceToDoctorUseCase handles assigning a service to a doctor
type AssignServiceToDoctorUseCase struct {
	doctorServiceRepo repository.DoctorServiceRepository
	serviceRepo       repository.ServiceRepository
	userRepo          repository.UserRepository
}

// NewAssignServiceToDoctorUseCase creates a new instance
func NewAssignServiceToDoctorUseCase(
	doctorServiceRepo repository.DoctorServiceRepository,
	serviceRepo repository.ServiceRepository,
	userRepo repository.UserRepository,
) *AssignServiceToDoctorUseCase {
	return &AssignServiceToDoctorUseCase{
		doctorServiceRepo: doctorServiceRepo,
		serviceRepo:       serviceRepo,
		userRepo:          userRepo,
	}
}

// Execute assigns a service to a doctor
func (uc *AssignServiceToDoctorUseCase) Execute(ctx context.Context, userID, serviceID string) error {
	// Validate inputs
	if userID == "" {
		return errors.New("doctor ID is required")
	}

	if serviceID == "" {
		return errors.New("service ID is required")
	}

	// Verify doctor exists and is active
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("doctor not found")
	}

	if user.Role != domain.RoleDoctor {
		return errors.New("user is not a doctor")
	}

	if !user.IsActive {
		return errors.New("doctor is not active")
	}

	// Get the actual doctor.id from doctors table
	doctorID, err := uc.userRepo.FindDoctorIDByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify service exists and is active
	service, err := uc.serviceRepo.FindByID(ctx, serviceID)
	if err != nil {
		return err
	}

	if service == nil {
		return errors.New("service not found")
	}

	if !service.IsActive {
		return errors.New("service is not active")
	}

	// Check if already assigned (using real doctor.id)
	isAssigned, err := uc.doctorServiceRepo.IsAssigned(ctx, doctorID, serviceID)
	if err != nil {
		return err
	}

	if isAssigned {
		return errors.New("service is already assigned to this doctor")
	}

	// Create assignment with the real doctor.id
	now := time.Now()
	doctorService := &domain.DoctorService{
		ID:        uuid.New().String(),
		DoctorID:  doctorID, // Now using the real doctor.id
		ServiceID: serviceID,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Validate domain entity
	if err := doctorService.Validate(); err != nil {
		return err
	}

	// Save to repository
	return uc.doctorServiceRepo.Assign(ctx, doctorService)
}

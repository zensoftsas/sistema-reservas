package service

import (
	"context"
	"errors"
	"strings"

	"version-1-0/internal/repository"
)

// DeleteServiceUseCase handles business logic for deleting a service
type DeleteServiceUseCase struct {
	serviceRepo       repository.ServiceRepository
	doctorServiceRepo repository.DoctorServiceRepository
}

// NewDeleteServiceUseCase creates a new instance
func NewDeleteServiceUseCase(
	serviceRepo repository.ServiceRepository,
	doctorServiceRepo repository.DoctorServiceRepository,
) *DeleteServiceUseCase {
	return &DeleteServiceUseCase{
		serviceRepo:       serviceRepo,
		doctorServiceRepo: doctorServiceRepo,
	}
}

// Execute deletes a service by ID
// Returns error if service has doctors assigned
func (uc *DeleteServiceUseCase) Execute(ctx context.Context, serviceID string) error {
	// Validate service ID
	if strings.TrimSpace(serviceID) == "" {
		return errors.New("service ID is required")
	}

	// Verify service exists
	service, err := uc.serviceRepo.FindByID(ctx, serviceID)
	if err != nil || service == nil {
		return errors.New("service not found")
	}

	// Check if service has doctors assigned
	doctors, err := uc.doctorServiceRepo.FindDoctorsByService(ctx, serviceID)
	if err != nil {
		return errors.New("failed to check service assignments")
	}

	if len(doctors) > 0 {
		return errors.New("cannot delete service with assigned doctors")
	}

	// Delete service
	if err := uc.serviceRepo.Delete(ctx, serviceID); err != nil {
		return errors.New("failed to delete service")
	}

	return nil
}

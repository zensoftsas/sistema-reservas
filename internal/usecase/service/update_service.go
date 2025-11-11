package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// UpdateServiceUseCase handles business logic for updating a service
type UpdateServiceUseCase struct {
	serviceRepo repository.ServiceRepository
}

// NewUpdateServiceUseCase creates a new instance
func NewUpdateServiceUseCase(serviceRepo repository.ServiceRepository) *UpdateServiceUseCase {
	return &UpdateServiceUseCase{
		serviceRepo: serviceRepo,
	}
}

// Execute updates a service by ID
func (uc *UpdateServiceUseCase) Execute(ctx context.Context, serviceID string, req UpdateServiceRequest) (*domain.Service, error) {
	// Validate service ID
	if strings.TrimSpace(serviceID) == "" {
		return nil, errors.New("service ID is required")
	}

	// Find existing service
	service, err := uc.serviceRepo.FindByID(ctx, serviceID)
	if err != nil || service == nil {
		return nil, errors.New("service not found")
	}

	// Update fields if provided
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		service.Name = *req.Name
	}

	if req.Description != nil {
		service.Description = *req.Description
	}

	if req.DurationMinutes != nil && *req.DurationMinutes > 0 {
		service.DurationMinutes = *req.DurationMinutes
	}

	if req.Price != nil && *req.Price >= 0 {
		service.Price = *req.Price
	}

	if req.IsActive != nil {
		service.IsActive = *req.IsActive
	}

	// Update timestamp
	service.UpdatedAt = time.Now()

	// Update in database
	if err := uc.serviceRepo.Update(ctx, service); err != nil {
		return nil, errors.New("failed to update service")
	}

	return service, nil
}

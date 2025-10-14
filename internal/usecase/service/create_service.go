package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// CreateServiceUseCase handles the creation of a new service
type CreateServiceUseCase struct {
	serviceRepo repository.ServiceRepository
}

// NewCreateServiceUseCase creates a new instance of CreateServiceUseCase
func NewCreateServiceUseCase(serviceRepo repository.ServiceRepository) *CreateServiceUseCase {
	return &CreateServiceUseCase{
		serviceRepo: serviceRepo,
	}
}

// Execute creates a new service
func (uc *CreateServiceUseCase) Execute(ctx context.Context, req CreateServiceRequest) (*CreateServiceResponse, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("service name is required")
	}

	if req.DurationMinutes <= 0 {
		return nil, errors.New("duration must be greater than 0")
	}

	if req.Price < 0 {
		return nil, errors.New("price cannot be negative")
	}

	// Create service entity
	now := time.Now()
	service := &domain.Service{
		ID:              uuid.New().String(),
		Name:            req.Name,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
		Price:           req.Price,
		IsActive:        true, // New services are active by default
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	// Validate domain entity
	if err := service.Validate(); err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.serviceRepo.Create(ctx, service); err != nil {
		return nil, err
	}

	// Return response
	return &CreateServiceResponse{
		ID:              service.ID,
		Name:            service.Name,
		Description:     service.Description,
		DurationMinutes: service.DurationMinutes,
		Price:           service.Price,
		IsActive:        service.IsActive,
		CreatedAt:       service.CreatedAt,
	}, nil
}

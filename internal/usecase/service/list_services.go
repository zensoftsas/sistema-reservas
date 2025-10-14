package service

import (
	"context"

	"version-1-0/internal/repository"
)

// ListServicesUseCase handles retrieving all active services
type ListServicesUseCase struct {
	serviceRepo repository.ServiceRepository
}

// NewListServicesUseCase creates a new instance of ListServicesUseCase
func NewListServicesUseCase(serviceRepo repository.ServiceRepository) *ListServicesUseCase {
	return &ListServicesUseCase{
		serviceRepo: serviceRepo,
	}
}

// Execute retrieves all active services
func (uc *ListServicesUseCase) Execute(ctx context.Context) ([]ServiceResponse, error) {
	// Retrieve services from repository
	services, err := uc.serviceRepo.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	responses := make([]ServiceResponse, len(services))
	for i, svc := range services {
		responses[i] = ServiceResponse{
			ID:              svc.ID,
			Name:            svc.Name,
			Description:     svc.Description,
			DurationMinutes: svc.DurationMinutes,
			Price:           svc.Price,
			IsActive:        svc.IsActive,
			CreatedAt:       svc.CreatedAt,
			UpdatedAt:       svc.UpdatedAt,
		}
	}

	return responses, nil
}

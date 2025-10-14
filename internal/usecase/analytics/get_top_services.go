package analytics

import (
	"context"

	"version-1-0/internal/repository"
)

// GetTopServicesUseCase handles retrieving top services statistics
type GetTopServicesUseCase struct {
	appointmentRepo repository.AppointmentRepository
}

// NewGetTopServicesUseCase creates a new instance
func NewGetTopServicesUseCase(
	appointmentRepo repository.AppointmentRepository,
) *GetTopServicesUseCase {
	return &GetTopServicesUseCase{
		appointmentRepo: appointmentRepo,
	}
}

// Execute retrieves most popular services
func (uc *GetTopServicesUseCase) Execute(ctx context.Context, limit int) ([]*TopService, error) {
	if limit <= 0 {
		limit = 10 // Default to top 10
	}

	// Get top services from repository
	servicesData, err := uc.appointmentRepo.GetTopServices(ctx, limit)
	if err != nil {
		return nil, err
	}

	var results []*TopService
	for _, data := range servicesData {
		results = append(results, &TopService{
			ServiceID:   data.ServiceID,
			ServiceName: data.ServiceName,
			TotalCitas:  data.Count,
		})
	}

	return results, nil
}

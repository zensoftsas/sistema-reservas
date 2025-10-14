package analytics

import (
	"context"

	"version-1-0/internal/repository"
)

// GetRevenueStatsUseCase handles retrieving revenue statistics
type GetRevenueStatsUseCase struct {
	appointmentRepo repository.AppointmentRepository
}

// NewGetRevenueStatsUseCase creates a new instance
func NewGetRevenueStatsUseCase(
	appointmentRepo repository.AppointmentRepository,
) *GetRevenueStatsUseCase {
	return &GetRevenueStatsUseCase{
		appointmentRepo: appointmentRepo,
	}
}

// Execute retrieves revenue statistics by service
func (uc *GetRevenueStatsUseCase) Execute(ctx context.Context) ([]*RevenueByService, error) {
	// Get revenue by service from repository
	revenueMap, err := uc.appointmentRepo.GetRevenueByService(ctx)
	if err != nil {
		return nil, err
	}

	// Convert to slice
	var results []*RevenueByService
	for serviceID, data := range revenueMap {
		results = append(results, &RevenueByService{
			ServiceID:   serviceID,
			ServiceName: data.ServiceName,
			TotalCitas:  data.Count,
			Revenue:     data.Revenue,
		})
	}

	return results, nil
}

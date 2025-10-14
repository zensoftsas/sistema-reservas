package analytics

import (
	"context"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// GetDashboardSummaryUseCase handles retrieving dashboard summary
type GetDashboardSummaryUseCase struct {
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository
}

// NewGetDashboardSummaryUseCase creates a new instance
func NewGetDashboardSummaryUseCase(
	appointmentRepo repository.AppointmentRepository,
	userRepo repository.UserRepository,
) *GetDashboardSummaryUseCase {
	return &GetDashboardSummaryUseCase{
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
	}
}

// Execute retrieves the dashboard summary statistics
func (uc *GetDashboardSummaryUseCase) Execute(ctx context.Context) (*DashboardSummary, error) {
	summary := &DashboardSummary{}

	// Count all appointments
	totalAppointments, err := uc.appointmentRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}
	summary.TotalAppointments = totalAppointments

	// Count appointments by status
	pending, err := uc.appointmentRepo.CountByStatus(ctx, string(domain.StatusPending))
	if err != nil {
		return nil, err
	}
	summary.PendingAppointments = pending

	confirmed, err := uc.appointmentRepo.CountByStatus(ctx, string(domain.StatusConfirmed))
	if err != nil {
		return nil, err
	}
	summary.ConfirmedAppointments = confirmed

	completed, err := uc.appointmentRepo.CountByStatus(ctx, string(domain.StatusCompleted))
	if err != nil {
		return nil, err
	}
	summary.CompletedAppointments = completed

	cancelled, err := uc.appointmentRepo.CountByStatus(ctx, string(domain.StatusCancelled))
	if err != nil {
		return nil, err
	}
	summary.CancelledAppointments = cancelled

	// Count patients and doctors
	patients, err := uc.userRepo.CountByRole(ctx, string(domain.RolePatient))
	if err != nil {
		return nil, err
	}
	summary.TotalPatients = patients

	doctors, err := uc.userRepo.CountByRole(ctx, string(domain.RoleDoctor))
	if err != nil {
		return nil, err
	}
	summary.TotalDoctors = doctors

	// Calculate total revenue
	revenue, err := uc.appointmentRepo.GetTotalRevenue(ctx)
	if err != nil {
		return nil, err
	}
	summary.TotalRevenue = revenue

	// Calculate cancellation rate
	if totalAppointments > 0 {
		summary.CancellationRate = float64(cancelled) / float64(totalAppointments) * 100
	} else {
		summary.CancellationRate = 0
	}

	return summary, nil
}

package schedule

import (
	"context"
	"errors"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// GetDoctorSchedulesUseCase handles retrieving doctor schedules
type GetDoctorSchedulesUseCase struct {
	scheduleRepo repository.ScheduleRepository
	userRepo     repository.UserRepository
}

// NewGetDoctorSchedulesUseCase creates a new instance
func NewGetDoctorSchedulesUseCase(
	scheduleRepo repository.ScheduleRepository,
	userRepo repository.UserRepository,
) *GetDoctorSchedulesUseCase {
	return &GetDoctorSchedulesUseCase{
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
	}
}

// Execute retrieves all schedules for a doctor
func (uc *GetDoctorSchedulesUseCase) Execute(ctx context.Context, userID string) ([]*domain.Schedule, error) {
	// Validate doctor exists
	doctor, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}
	if doctor.Role != domain.RoleDoctor {
		return nil, errors.New("user is not a doctor")
	}

	// Get real doctor.id
	doctorID, err := uc.userRepo.FindDoctorIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get all schedules
	schedules, err := uc.scheduleRepo.FindByDoctor(ctx, doctorID)
	if err != nil {
		return nil, err
	}

	return schedules, nil
}

package schedule

import (
	"context"
	"errors"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// UpdateScheduleUseCase handles updating doctor schedules
type UpdateScheduleUseCase struct {
	scheduleRepo repository.ScheduleRepository
	userRepo     repository.UserRepository
}

// NewUpdateScheduleUseCase creates a new instance
func NewUpdateScheduleUseCase(
	scheduleRepo repository.ScheduleRepository,
	userRepo repository.UserRepository,
) *UpdateScheduleUseCase {
	return &UpdateScheduleUseCase{
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
	}
}

// Execute updates an existing schedule
func (uc *UpdateScheduleUseCase) Execute(
	ctx context.Context,
	scheduleID string,
	userID string,
	dayOfWeek string,
	startTime string,
	endTime string,
	slotDuration int,
	isActive bool,
) (*domain.Schedule, error) {
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

	// Find existing schedule
	existingSchedule, err := uc.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return nil, err
	}
	if existingSchedule == nil {
		return nil, errors.New("schedule not found")
	}

	// Verify ownership
	if existingSchedule.DoctorID != doctorID {
		return nil, errors.New("schedule does not belong to this doctor")
	}

	// Update fields
	existingSchedule.DayOfWeek = dayOfWeek
	existingSchedule.StartTime = startTime
	existingSchedule.EndTime = endTime
	existingSchedule.SlotDuration = slotDuration
	existingSchedule.IsActive = isActive
	existingSchedule.UpdatedAt = time.Now()

	// Validate
	if err := existingSchedule.Validate(); err != nil {
		return nil, err
	}

	// Check for overlapping schedules (excluding this one)
	existing, err := uc.scheduleRepo.FindByDoctorAndDay(ctx, doctorID, dayOfWeek)
	if err != nil {
		return nil, err
	}

	for _, sched := range existing {
		if sched.ID != scheduleID && schedulesOverlap(existingSchedule, sched) {
			return nil, errors.New("schedule overlaps with existing schedule for this day")
		}
	}

	// Save
	if err := uc.scheduleRepo.Update(ctx, existingSchedule); err != nil {
		return nil, err
	}

	return existingSchedule, nil
}

package schedule

import (
	"context"
	"errors"

	"version-1-0/internal/repository"
)

// DeleteScheduleUseCase handles deleting schedules
type DeleteScheduleUseCase struct {
	scheduleRepo repository.ScheduleRepository
}

// NewDeleteScheduleUseCase creates a new instance
func NewDeleteScheduleUseCase(
	scheduleRepo repository.ScheduleRepository,
) *DeleteScheduleUseCase {
	return &DeleteScheduleUseCase{
		scheduleRepo: scheduleRepo,
	}
}

// Execute deletes a schedule by ID
func (uc *DeleteScheduleUseCase) Execute(ctx context.Context, scheduleID string) error {
	// Verify schedule exists
	schedule, err := uc.scheduleRepo.FindByID(ctx, scheduleID)
	if err != nil {
		return err
	}
	if schedule == nil {
		return errors.New("schedule not found")
	}

	// Delete
	return uc.scheduleRepo.Delete(ctx, scheduleID)
}

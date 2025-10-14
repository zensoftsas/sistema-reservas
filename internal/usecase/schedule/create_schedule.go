package schedule

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// CreateScheduleUseCase handles creating doctor schedules
type CreateScheduleUseCase struct {
	scheduleRepo repository.ScheduleRepository
	userRepo     repository.UserRepository
}

// NewCreateScheduleUseCase creates a new instance
func NewCreateScheduleUseCase(
	scheduleRepo repository.ScheduleRepository,
	userRepo repository.UserRepository,
) *CreateScheduleUseCase {
	return &CreateScheduleUseCase{
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
	}
}

// Execute creates a new schedule for a doctor
func (uc *CreateScheduleUseCase) Execute(
	ctx context.Context,
	userID string,
	dayOfWeek string,
	startTime string,
	endTime string,
	slotDuration int,
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

	// Create schedule entity
	now := time.Now()
	schedule := &domain.Schedule{
		ID:           uuid.New().String(),
		DoctorID:     doctorID,
		DayOfWeek:    dayOfWeek,
		StartTime:    startTime,
		EndTime:      endTime,
		SlotDuration: slotDuration,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Validate
	if err := schedule.Validate(); err != nil {
		return nil, err
	}

	// Check for overlapping schedules
	existing, err := uc.scheduleRepo.FindByDoctorAndDay(ctx, doctorID, dayOfWeek)
	if err != nil {
		return nil, err
	}

	for _, existingSched := range existing {
		if schedulesOverlap(schedule, existingSched) {
			return nil, errors.New("schedule overlaps with existing schedule for this day")
		}
	}

	// Save
	if err := uc.scheduleRepo.Create(ctx, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

// schedulesOverlap checks if two schedules overlap
func schedulesOverlap(s1, s2 *domain.Schedule) bool {
	start1, _ := time.Parse("15:04", s1.StartTime)
	end1, _ := time.Parse("15:04", s1.EndTime)
	start2, _ := time.Parse("15:04", s2.StartTime)
	end2, _ := time.Parse("15:04", s2.EndTime)

	return start1.Before(end2) && end1.After(start2)
}

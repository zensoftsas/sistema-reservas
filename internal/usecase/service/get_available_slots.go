package service

import (
	"context"
	"errors"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// TimeSlot represents an available time slot
type TimeSlot struct {
	Time      string `json:"time"`       // HH:MM format
	Available bool   `json:"available"`
}

// GetAvailableSlotsUseCase calculates available time slots for a doctor-service combination
type GetAvailableSlotsUseCase struct {
	serviceRepo     repository.ServiceRepository
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository
	scheduleRepo    repository.ScheduleRepository
}

// NewGetAvailableSlotsUseCase creates a new instance
func NewGetAvailableSlotsUseCase(
	serviceRepo repository.ServiceRepository,
	appointmentRepo repository.AppointmentRepository,
	userRepo repository.UserRepository,
	scheduleRepo repository.ScheduleRepository,
) *GetAvailableSlotsUseCase {
	return &GetAvailableSlotsUseCase{
		serviceRepo:     serviceRepo,
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
		scheduleRepo:    scheduleRepo,
	}
}

// Execute calculates available slots for a doctor offering a service on a specific date
func (uc *GetAvailableSlotsUseCase) Execute(ctx context.Context, userID, serviceID string, date time.Time) ([]TimeSlot, error) {
	// Validate inputs
	if userID == "" {
		return nil, errors.New("doctor ID is required")
	}
	if serviceID == "" {
		return nil, errors.New("service ID is required")
	}

	// Verify doctor exists
	doctor, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if doctor == nil {
		return nil, errors.New("doctor not found")
	}

	// Get doctor.id real
	doctorID, err := uc.userRepo.FindDoctorIDByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Verify service exists
	service, err := uc.serviceRepo.FindByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, errors.New("service not found")
	}

	// Get doctor's schedule for this day of week
	dayOfWeek := domain.GetDayOfWeekFromDate(date)
	schedules, err := uc.scheduleRepo.FindByDoctorAndDay(ctx, doctorID, dayOfWeek)
	if err != nil {
		return nil, err
	}

	// If doctor doesn't work this day, return empty slots
	if len(schedules) == 0 {
		return []TimeSlot{}, nil
	}

	// Generate slots for each schedule block
	var allSlots []TimeSlot
	for _, sched := range schedules {
		// Parse start and end times
		start, _ := time.Parse("15:04", sched.StartTime)
		end, _ := time.Parse("15:04", sched.EndTime)

		startHour := start.Hour()
		startMinute := start.Minute()
		endHour := end.Hour()
		endMinute := end.Minute()

		startMinutes := startHour*60 + startMinute
		endMinutes := endHour*60 + endMinute

		// Generate slots for this schedule block
		slots := generateTimeSlotsFromMinutes(startMinutes, endMinutes, service.DurationMinutes)
		allSlots = append(allSlots, slots...)
	}

	slots := allSlots

	// Get existing appointments for this doctor on this date
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	appointments, err := uc.appointmentRepo.FindByDoctorAndDateRange(ctx, doctorID, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	// Mark slots as unavailable if they conflict with existing appointments
	for i := range slots {
		slots[i].Available = true
		slotTime := parseTimeSlot(date, slots[i].Time)

		for _, apt := range appointments {
			if apt.Status == "cancelled" {
				continue
			}

			// Check if slot overlaps with appointment
			aptEnd := apt.ScheduledAt.Add(time.Duration(apt.Duration) * time.Minute)
			slotEnd := slotTime.Add(time.Duration(service.DurationMinutes) * time.Minute)

			if slotTime.Before(aptEnd) && slotEnd.After(apt.ScheduledAt) {
				slots[i].Available = false
				break
			}
		}
	}

	return slots, nil
}

// generateTimeSlots creates time slots from start to end hour with given duration
func generateTimeSlots(startHour, endHour, durationMinutes int) []TimeSlot {
	var slots []TimeSlot

	currentMinutes := startHour * 60
	endMinutes := endHour * 60

	for currentMinutes < endMinutes {
		hours := currentMinutes / 60
		minutes := currentMinutes % 60
		timeStr := formatTime(hours, minutes)

		slots = append(slots, TimeSlot{
			Time:      timeStr,
			Available: false,
		})

		currentMinutes += durationMinutes
	}

	return slots
}

// formatTime formats hours and minutes as HH:MM
func formatTime(hours, minutes int) string {
	return time.Date(0, 1, 1, hours, minutes, 0, 0, time.UTC).Format("15:04")
}

// parseTimeSlot parses a time string (HH:MM) and combines with date
func parseTimeSlot(date time.Time, timeStr string) time.Time {
	t, _ := time.Parse("15:04", timeStr)
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, date.Location())
}

// generateTimeSlotsFromMinutes creates time slots from start to end minutes with given duration
func generateTimeSlotsFromMinutes(startMinutes, endMinutes, durationMinutes int) []TimeSlot {
	var slots []TimeSlot

	currentMinutes := startMinutes

	for currentMinutes < endMinutes {
		hours := currentMinutes / 60
		minutes := currentMinutes % 60
		timeStr := formatTime(hours, minutes)

		slots = append(slots, TimeSlot{
			Time:      timeStr,
			Available: false,
		})

		currentMinutes += durationMinutes
	}

	return slots
}

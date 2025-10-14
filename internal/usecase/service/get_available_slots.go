package service

import (
	"context"
	"errors"
	"time"

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
}

// NewGetAvailableSlotsUseCase creates a new instance
func NewGetAvailableSlotsUseCase(
	serviceRepo repository.ServiceRepository,
	appointmentRepo repository.AppointmentRepository,
	userRepo repository.UserRepository,
) *GetAvailableSlotsUseCase {
	return &GetAvailableSlotsUseCase{
		serviceRepo:     serviceRepo,
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
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

	// TODO: Get doctor's schedule for this day of week
	// For now, use a default schedule: 9:00 AM to 5:00 PM
	startHour := 9
	endHour := 17

	// Generate time slots based on service duration
	slots := generateTimeSlots(startHour, endHour, service.DurationMinutes)

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

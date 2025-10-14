package domain

import (
	"errors"
	"time"
)

// Schedule represents a doctor's working schedule for a specific day
type Schedule struct {
	ID           string    `json:"id"`
	DoctorID     string    `json:"doctor_id"`
	DayOfWeek    string    `json:"day_of_week"`    // monday, tuesday, etc.
	StartTime    string    `json:"start_time"`     // HH:MM format
	EndTime      string    `json:"end_time"`       // HH:MM format
	SlotDuration int       `json:"slot_duration"`  // Minutes per slot
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// DayOfWeek constants
const (
	Monday    = "monday"
	Tuesday   = "tuesday"
	Wednesday = "wednesday"
	Thursday  = "thursday"
	Friday    = "friday"
	Saturday  = "saturday"
	Sunday    = "sunday"
)

// ValidDaysOfWeek contains all valid day names
var ValidDaysOfWeek = []string{
	Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday,
}

// Validate validates the schedule data
func (s *Schedule) Validate() error {
	if s.DoctorID == "" {
		return errors.New("doctor_id is required")
	}

	// Validate day of week
	validDay := false
	for _, day := range ValidDaysOfWeek {
		if s.DayOfWeek == day {
			validDay = true
			break
		}
	}
	if !validDay {
		return errors.New("invalid day_of_week, must be monday-sunday")
	}

	// Validate time format
	if !isValidTimeFormat(s.StartTime) {
		return errors.New("start_time must be in HH:MM format")
	}
	if !isValidTimeFormat(s.EndTime) {
		return errors.New("end_time must be in HH:MM format")
	}

	// Validate start < end
	start, _ := time.Parse("15:04", s.StartTime)
	end, _ := time.Parse("15:04", s.EndTime)
	if !start.Before(end) {
		return errors.New("start_time must be before end_time")
	}

	if s.SlotDuration <= 0 {
		return errors.New("slot_duration must be greater than 0")
	}

	return nil
}

// Activate activates the schedule
func (s *Schedule) Activate() {
	s.IsActive = true
	s.UpdatedAt = time.Now()
}

// Deactivate deactivates the schedule
func (s *Schedule) Deactivate() {
	s.IsActive = false
	s.UpdatedAt = time.Now()
}

// isValidTimeFormat checks if a time string is in HH:MM format
func isValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04", timeStr)
	return err == nil
}

// GetDayOfWeekFromDate returns the day of week name for a given date
func GetDayOfWeekFromDate(date time.Time) string {
	switch date.Weekday() {
	case time.Monday:
		return Monday
	case time.Tuesday:
		return Tuesday
	case time.Wednesday:
		return Wednesday
	case time.Thursday:
		return Thursday
	case time.Friday:
		return Friday
	case time.Saturday:
		return Saturday
	case time.Sunday:
		return Sunday
	default:
		return ""
	}
}

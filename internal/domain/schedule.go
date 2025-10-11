package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// DayOfWeek represents a day of the week
type DayOfWeek string

// Day of week constants
const (
	Monday    DayOfWeek = "monday"
	Tuesday   DayOfWeek = "tuesday"
	Wednesday DayOfWeek = "wednesday"
	Thursday  DayOfWeek = "thursday"
	Friday    DayOfWeek = "friday"
	Saturday  DayOfWeek = "saturday"
	Sunday    DayOfWeek = "sunday"
)

// Schedule represents a doctor's schedule configuration in the medical reservation system
type Schedule struct {
	ID           string    `json:"id"`
	DoctorID     string    `json:"doctor_id"`
	DayOfWeek    DayOfWeek `json:"day_of_week"`
	StartTime    string    `json:"start_time"` // Format: "HH:MM"
	EndTime      string    `json:"end_time"`   // Format: "HH:MM"
	SlotDuration int       `json:"slot_duration"` // in minutes
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Validate checks if the Schedule entity has all required fields properly set
// Returns an error if any validation rule fails
func (s *Schedule) Validate() error {
	if strings.TrimSpace(s.ID) == "" {
		return errors.New("schedule ID is required")
	}

	if strings.TrimSpace(s.DoctorID) == "" {
		return errors.New("schedule doctor ID is required")
	}

	if !IsValidDayOfWeek(string(s.DayOfWeek)) {
		return errors.New("invalid day of week")
	}

	if !IsValidTimeFormat(s.StartTime) {
		return errors.New("invalid start time format, must be HH:MM")
	}

	if !IsValidTimeFormat(s.EndTime) {
		return errors.New("invalid end time format, must be HH:MM")
	}

	// Compare start time and end time
	startMinutes, _ := parseTimeToMinutes(s.StartTime)
	endMinutes, _ := parseTimeToMinutes(s.EndTime)

	if startMinutes >= endMinutes {
		return errors.New("start time must be before end time")
	}

	if s.SlotDuration <= 0 {
		return errors.New("slot duration must be greater than 0")
	}

	if s.CreatedAt.IsZero() {
		return errors.New("schedule created at is required")
	}

	if s.UpdatedAt.IsZero() {
		return errors.New("schedule updated at is required")
	}

	return nil
}

// TotalSlots calculates how many appointment slots fit between StartTime and EndTime
// based on the SlotDuration
func (s *Schedule) TotalSlots() int {
	startMinutes, err := parseTimeToMinutes(s.StartTime)
	if err != nil {
		return 0
	}

	endMinutes, err := parseTimeToMinutes(s.EndTime)
	if err != nil {
		return 0
	}

	totalMinutes := endMinutes - startMinutes
	if totalMinutes <= 0 {
		return 0
	}

	return totalMinutes / s.SlotDuration
}

// IsValidTimeFormat verifies if a time string is in valid "HH:MM" format
// Hours must be 00-23, minutes must be 00-59
func IsValidTimeFormat(timeStr string) bool {
	// Regular expression to match HH:MM format
	pattern := `^([0-1][0-9]|2[0-3]):([0-5][0-9])$`
	matched, _ := regexp.MatchString(pattern, timeStr)
	return matched
}

// IsValidDayOfWeek checks if a given string is a valid DayOfWeek
func IsValidDayOfWeek(day string) bool {
	d := DayOfWeek(strings.ToLower(day))
	return d == Monday || d == Tuesday || d == Wednesday || d == Thursday ||
		d == Friday || d == Saturday || d == Sunday
}

// parseTimeToMinutes converts a time string in "HH:MM" format to total minutes since midnight
// Returns the number of minutes and an error if the format is invalid
func parseTimeToMinutes(timeStr string) (int, error) {
	if !IsValidTimeFormat(timeStr) {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}

	parts := strings.Split(timeStr, ":")
	hours, _ := strconv.Atoi(parts[0])
	minutes, _ := strconv.Atoi(parts[1])

	return hours*60 + minutes, nil
}

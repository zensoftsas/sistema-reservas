package reminder

import (
	"context"
	"log"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
	"version-1-0/pkg/email"
)

// ReminderService handles sending appointment reminders
type ReminderService struct {
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository
	emailService    *email.EmailService
}

// NewReminderService creates a new reminder service
func NewReminderService(
	appointmentRepo repository.AppointmentRepository,
	userRepo repository.UserRepository,
	emailService *email.EmailService,
) *ReminderService {
	return &ReminderService{
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
		emailService:    emailService,
	}
}

// Start begins the reminder scheduler
// Runs every 10 minutes checking for appointments that need reminders
func (s *ReminderService) Start() {
	log.Println("Reminder service started - checking every 10 minutes")

	// Run immediately on start
	s.checkAndSendReminders()

	// Then run every 10 minutes
	ticker := time.NewTicker(10 * time.Minute)

	go func() {
		for range ticker.C {
			s.checkAndSendReminders()
		}
	}()
}

// checkAndSendReminders checks for appointments needing reminders and sends them
func (s *ReminderService) checkAndSendReminders() {
	ctx := context.Background()
	now := time.Now()

	log.Println("Checking for appointments needing reminders...")

	// Check 24-hour reminders
	s.send24HourReminders(ctx, now)

	// Check 1-hour reminders (optional, uncomment to enable)
	// s.send1HourReminders(ctx, now)
}

// send24HourReminders sends reminders for appointments 24 hours from now
func (s *ReminderService) send24HourReminders(ctx context.Context, now time.Time) {
	// Calculate time window: 24 hours from now (+/- 10 minutes buffer)
	targetTime := now.Add(24 * time.Hour)
	startWindow := targetTime.Add(-10 * time.Minute)
	endWindow := targetTime.Add(10 * time.Minute)

	// Find all confirmed appointments in the time window
	appointments, err := s.findAppointmentsInWindow(ctx, startWindow, endWindow, "confirmed")
	if err != nil {
		log.Printf("Error finding appointments for 24h reminders: %v", err)
		return
	}

	sent := 0
	for _, apt := range appointments {
		// Skip if reminder already sent
		if apt.Reminder24hSent {
			continue
		}

		// Get patient and doctor info
		patient, _ := s.userRepo.FindByID(ctx, apt.PatientID)
		doctor, _ := s.userRepo.FindByID(ctx, apt.DoctorID)

		if patient == nil || doctor == nil {
			continue
		}

		// Send reminder email
		err := s.send24HourReminderEmail(patient, doctor, apt)
		if err != nil {
			log.Printf("Error sending 24h reminder for appointment %s: %v", apt.ID, err)
			continue
		}

		// Mark as sent
		s.markReminder24hSent(ctx, apt.ID)
		sent++
	}

	if sent > 0 {
		log.Printf("Sent %d 24-hour reminders", sent)
	}
}

// send24HourReminderEmail sends the 24-hour reminder email
func (s *ReminderService) send24HourReminderEmail(patient, doctor *domain.User, apt *domain.Appointment) error {
	if s.emailService == nil {
		return nil
	}

	patientName := patient.FirstName + " " + patient.LastName
	doctorName := doctor.FirstName + " " + doctor.LastName
	date := apt.ScheduledAt.Format("2006-01-02")
	timeStr := apt.ScheduledAt.Format("15:04")

	return s.emailService.SendAppointmentReminder(
		patient.Email,
		patientName,
		doctorName,
		date,
		timeStr,
		"24 horas",
	)
}

// findAppointmentsInWindow finds appointments in a time window
func (s *ReminderService) findAppointmentsInWindow(ctx context.Context, start, end time.Time, status string) ([]*domain.Appointment, error) {
	return s.appointmentRepo.FindByScheduledAtRange(ctx, start, end, status)
}

// markReminder24hSent marks an appointment's 24h reminder as sent
func (s *ReminderService) markReminder24hSent(ctx context.Context, appointmentID string) {
	err := s.appointmentRepo.MarkReminder24hSent(ctx, appointmentID)
	if err != nil {
		log.Printf("Error marking 24h reminder sent for %s: %v", appointmentID, err)
	}
}

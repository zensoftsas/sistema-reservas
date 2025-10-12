package repository

import (
	"context"
	"time"

	"version-1-0/internal/domain"
)

// UserRepository defines the interface for user data persistence operations
type UserRepository interface {
	// Create inserts a new user into the repository
	Create(ctx context.Context, user *domain.User) error

	// FindByID retrieves a user by their unique identifier
	FindByID(ctx context.Context, id string) (*domain.User, error)

	// FindByEmail retrieves a user by their email address
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update modifies an existing user in the repository
	Update(ctx context.Context, user *domain.User) error

	// Delete removes a user from the repository by their ID
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of users
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)

	// FindDoctorsBySpecialty retrieves all active doctors filtered by specialty
	FindDoctorsBySpecialty(ctx context.Context, specialty string) ([]*domain.User, error)

	// GetAllDoctors retrieves all active doctors
	GetAllDoctors(ctx context.Context) ([]*domain.User, error)
}

// PatientRepository defines the interface for patient data persistence operations
type PatientRepository interface {
	// Create inserts a new patient into the repository
	Create(ctx context.Context, patient *domain.Patient) error

	// FindByID retrieves a patient by their unique identifier
	FindByID(ctx context.Context, id string) (*domain.Patient, error)

	// FindByUserID retrieves a patient by their associated user ID
	FindByUserID(ctx context.Context, userID string) (*domain.Patient, error)

	// Update modifies an existing patient in the repository
	Update(ctx context.Context, patient *domain.Patient) error

	// Delete removes a patient from the repository by their ID
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of patients
	List(ctx context.Context, limit, offset int) ([]*domain.Patient, error)
}

// DoctorRepository defines the interface for doctor data persistence operations
type DoctorRepository interface {
	// Create inserts a new doctor into the repository
	Create(ctx context.Context, doctor *domain.Doctor) error

	// FindByID retrieves a doctor by their unique identifier
	FindByID(ctx context.Context, id string) (*domain.Doctor, error)

	// FindByUserID retrieves a doctor by their associated user ID
	FindByUserID(ctx context.Context, userID string) (*domain.Doctor, error)

	// FindBySpecialty retrieves all doctors with a specific specialty
	FindBySpecialty(ctx context.Context, specialty string) ([]*domain.Doctor, error)

	// Update modifies an existing doctor in the repository
	Update(ctx context.Context, doctor *domain.Doctor) error

	// Delete removes a doctor from the repository by their ID
	Delete(ctx context.Context, id string) error

	// List retrieves a paginated list of doctors
	List(ctx context.Context, limit, offset int) ([]*domain.Doctor, error)
}

// AppointmentRepository defines the interface for appointment data persistence operations
type AppointmentRepository interface {
	// Create inserts a new appointment into the repository
	Create(ctx context.Context, appointment *domain.Appointment) error

	// FindByID retrieves an appointment by its unique identifier
	FindByID(ctx context.Context, id string) (*domain.Appointment, error)

	// FindByPatientID retrieves all appointments for a specific patient
	FindByPatientID(ctx context.Context, patientID string) ([]*domain.Appointment, error)

	// FindByDoctorID retrieves all appointments for a specific doctor
	FindByDoctorID(ctx context.Context, doctorID string) ([]*domain.Appointment, error)

	// FindByDoctorAndDate retrieves all appointments for a doctor on a specific date
	FindByDoctorAndDate(ctx context.Context, doctorID string, date time.Time) ([]*domain.Appointment, error)

	// Update modifies an existing appointment in the repository
	Update(ctx context.Context, appointment *domain.Appointment) error

	// Delete removes an appointment from the repository by its ID
	Delete(ctx context.Context, id string) error

	// FindByScheduledAtRange retrieves appointments within a time range with specific status
	FindByScheduledAtRange(ctx context.Context, start, end time.Time, status string) ([]*domain.Appointment, error)

	// MarkReminder24hSent marks the 24-hour reminder as sent for an appointment
	MarkReminder24hSent(ctx context.Context, id string) error

	// MarkReminder1hSent marks the 1-hour reminder as sent for an appointment
	MarkReminder1hSent(ctx context.Context, id string) error
}

// ScheduleRepository defines the interface for schedule data persistence operations
type ScheduleRepository interface {
	// Create inserts a new schedule into the repository
	Create(ctx context.Context, schedule *domain.Schedule) error

	// FindByID retrieves a schedule by its unique identifier
	FindByID(ctx context.Context, id string) (*domain.Schedule, error)

	// FindByDoctorID retrieves all schedules for a specific doctor
	FindByDoctorID(ctx context.Context, doctorID string) ([]*domain.Schedule, error)

	// FindByDoctorAndDay retrieves a doctor's schedule for a specific day of the week
	FindByDoctorAndDay(ctx context.Context, doctorID string, day domain.DayOfWeek) (*domain.Schedule, error)

	// Update modifies an existing schedule in the repository
	Update(ctx context.Context, schedule *domain.Schedule) error

	// Delete removes a schedule from the repository by its ID
	Delete(ctx context.Context, id string) error
}

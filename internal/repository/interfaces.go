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

	// FindDoctorIDByUserID returns the doctor.id for a given user_id
	FindDoctorIDByUserID(ctx context.Context, userID string) (string, error)

	// Analytics methods
	// CountByRole counts users by role
	CountByRole(ctx context.Context, role string) (int, error)

	// CountAllActive counts all active users
	CountAllActive(ctx context.Context) (int, error)
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

	// FindByDoctorAndDateRange retrieves appointments for a doctor within a date range
	FindByDoctorAndDateRange(ctx context.Context, doctorID string, start, end time.Time) ([]*domain.Appointment, error)

	// MarkReminder24hSent marks the 24-hour reminder as sent for an appointment
	MarkReminder24hSent(ctx context.Context, id string) error

	// MarkReminder1hSent marks the 1-hour reminder as sent for an appointment
	MarkReminder1hSent(ctx context.Context, id string) error

	// Analytics methods
	// CountByStatus counts appointments by status
	CountByStatus(ctx context.Context, status string) (int, error)

	// CountAll counts all appointments
	CountAll(ctx context.Context) (int, error)

	// GetTotalRevenue calculates total revenue from completed appointments
	GetTotalRevenue(ctx context.Context) (float64, error)

	// GetRevenueByService gets revenue grouped by service
	GetRevenueByService(ctx context.Context) (map[string]struct {
		ServiceName string
		Count       int
		Revenue     float64
	}, error)

	// GetTopDoctors gets doctors with most appointments
	GetTopDoctors(ctx context.Context, limit int) ([]struct {
		DoctorID              string
		TotalAppointments     int
		CompletedAppointments int
	}, error)

	// GetTopServices gets most used services
	GetTopServices(ctx context.Context, limit int) ([]struct {
		ServiceID   string
		ServiceName string
		Count       int
	}, error)
}

// ScheduleRepository defines methods for schedule data access
type ScheduleRepository interface {
	// Create creates a new schedule
	Create(ctx context.Context, schedule *domain.Schedule) error

	// FindByID finds a schedule by ID
	FindByID(ctx context.Context, id string) (*domain.Schedule, error)

	// FindByDoctorAndDay finds schedules for a doctor on a specific day
	FindByDoctorAndDay(ctx context.Context, doctorID, dayOfWeek string) ([]*domain.Schedule, error)

	// FindByDoctor finds all schedules for a doctor
	FindByDoctor(ctx context.Context, doctorID string) ([]*domain.Schedule, error)

	// Update updates a schedule
	Update(ctx context.Context, schedule *domain.Schedule) error

	// Delete deletes a schedule
	Delete(ctx context.Context, id string) error

	// DeleteByDoctorAndDay deletes all schedules for a doctor on a specific day
	DeleteByDoctorAndDay(ctx context.Context, doctorID, dayOfWeek string) error
}

// ServiceRepository defines the interface for service data persistence operations
type ServiceRepository interface {
	// Create inserts a new service into the repository
	Create(ctx context.Context, service *domain.Service) error

	// FindByID retrieves a service by its unique identifier
	FindByID(ctx context.Context, id string) (*domain.Service, error)

	// ListActive retrieves all active services
	ListActive(ctx context.Context) ([]*domain.Service, error)

	// ListAll retrieves all services (active and inactive)
	ListAll(ctx context.Context) ([]*domain.Service, error)

	// Update modifies an existing service in the repository
	Update(ctx context.Context, service *domain.Service) error

	// Delete removes a service from the repository by its ID
	Delete(ctx context.Context, id string) error
}

// DoctorServiceRepository defines the interface for doctor-service relationship data persistence operations
type DoctorServiceRepository interface {
	// Assign creates a relationship between a doctor and a service
	Assign(ctx context.Context, doctorService *domain.DoctorService) error

	// Remove removes a service assignment from a doctor
	Remove(ctx context.Context, doctorID, serviceID string) error

	// FindDoctorsByService returns all doctors that offer a specific service
	FindDoctorsByService(ctx context.Context, serviceID string) ([]*domain.User, error)

	// FindServicesByDoctor returns all services offered by a specific doctor
	FindServicesByDoctor(ctx context.Context, doctorID string) ([]*domain.Service, error)

	// IsAssigned checks if a doctor is assigned to a service
	IsAssigned(ctx context.Context, doctorID, serviceID string) (bool, error)

	// FindByDoctorAndService retrieves a specific doctor-service relationship
	FindByDoctorAndService(ctx context.Context, doctorID, serviceID string) (*domain.DoctorService, error)
}

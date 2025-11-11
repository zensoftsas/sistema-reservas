package dto

// Authentication DTOs
type LoginRequest struct {
	Email    string `json:"email" example:"admin@clinica.com"`
	Password string `json:"password" example:"admin123"`
}

type LoginResponse struct {
	Token     string       `json:"token"`
	ExpiresAt string       `json:"expires_at"`
	User      UserResponse `json:"user"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
}

// User DTOs
type CreateUserRequest struct {
	Email     string `json:"email" example:"doctor@clinica.com"`
	Password  string `json:"password" example:"password123"`
	FirstName string `json:"first_name" example:"Juan"`
	LastName  string `json:"last_name" example:"Pérez"`
	Phone     string `json:"phone,omitempty" example:"+593999999999"`
	Role      string `json:"role" example:"doctor"`
}

// Appointment DTOs
type CreateAppointmentRequest struct {
	DoctorID        string `json:"doctor_id" example:"uuid"`
	ServiceID       string `json:"service_id" example:"uuid"`
	AppointmentDate string `json:"appointment_date" example:"2025-11-15"`
	AppointmentTime string `json:"appointment_time" example:"10:00"`
	Reason          string `json:"reason" example:"Consulta general"`
}

type AppointmentResponse struct {
	ID              string `json:"id"`
	PatientID       string `json:"patient_id"`
	DoctorID        string `json:"doctor_id"`
	ServiceID       string `json:"service_id"`
	ServiceName     string `json:"service_name"`
	ScheduledAt     string `json:"scheduled_at"`
	Duration        int    `json:"duration"`
	Reason          string `json:"reason"`
	Notes           string `json:"notes"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// Service DTOs
type CreateServiceRequest struct {
	Name            string  `json:"name" example:"Consulta General"`
	Description     string  `json:"description" example:"Consulta médica general"`
	DurationMinutes int     `json:"duration_minutes" example:"30"`
	Price           float64 `json:"price" example:"80.00"`
}

type ServiceResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	DurationMinutes int     `json:"duration_minutes"`
	Price           float64 `json:"price"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
}

type TimeSlot struct {
	StartTime string `json:"start_time" example:"09:00"`
	EndTime   string `json:"end_time" example:"09:30"`
}

// Schedule DTOs
type CreateScheduleRequest struct {
	DoctorID  string      `json:"doctor_id" example:"uuid"`
	DayOfWeek int         `json:"day_of_week" example:"1"`
	Blocks    []TimeBlock `json:"blocks"`
}

type TimeBlock struct {
	StartTime string `json:"start_time" example:"09:00"`
	EndTime   string `json:"end_time" example:"13:00"`
}

type ScheduleResponse struct {
	ID        string      `json:"id"`
	DoctorID  string      `json:"doctor_id"`
	DayOfWeek int         `json:"day_of_week"`
	Blocks    []TimeBlock `json:"blocks"`
	CreatedAt string      `json:"created_at"`
}

// Analytics DTOs
type DashboardSummary struct {
	TotalAppointments     int     `json:"total_appointments"`
	PendingAppointments   int     `json:"pending_appointments"`
	ConfirmedAppointments int     `json:"confirmed_appointments"`
	CompletedAppointments int     `json:"completed_appointments"`
	CancelledAppointments int     `json:"cancelled_appointments"`
	TotalPatients         int     `json:"total_patients"`
	TotalDoctors          int     `json:"total_doctors"`
	TotalRevenue          float64 `json:"total_revenue"`
	CancellationRate      float64 `json:"cancellation_rate"`
}

type RevenueByService struct {
	ServiceID   string  `json:"service_id"`
	ServiceName string  `json:"service_name"`
	TotalCitas  int     `json:"total_citas"`
	Revenue     float64 `json:"revenue"`
}

type TopDoctor struct {
	DoctorID              string `json:"doctor_id"`
	DoctorName            string `json:"doctor_name"`
	TotalAppointments     int    `json:"total_appointments"`
	CompletedAppointments int    `json:"completed_appointments"`
}

type TopService struct {
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	TotalCitas  int    `json:"total_citas"`
}

// Error response
type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

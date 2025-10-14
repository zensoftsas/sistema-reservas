package analytics

import "time"

// DashboardSummary represents the main dashboard summary
type DashboardSummary struct {
	TotalAppointments     int     `json:"total_appointments"`
	PendingAppointments   int     `json:"pending_appointments"`
	ConfirmedAppointments int     `json:"confirmed_appointments"`
	CompletedAppointments int     `json:"completed_appointments"`
	CancelledAppointments int     `json:"cancelled_appointments"`
	TotalPatients         int     `json:"total_patients"`
	TotalDoctors          int     `json:"total_doctors"`
	TotalRevenue          float64 `json:"total_revenue"`
	CancellationRate      float64 `json:"cancellation_rate"` // Percentage
}

// AppointmentsByPeriod represents appointments grouped by time period
type AppointmentsByPeriod struct {
	Period string `json:"period"` // Date or week
	Count  int    `json:"count"`
}

// RevenueByService represents revenue statistics per service
type RevenueByService struct {
	ServiceID   string  `json:"service_id"`
	ServiceName string  `json:"service_name"`
	TotalCitas  int     `json:"total_citas"`
	Revenue     float64 `json:"revenue"`
}

// TopDoctor represents a doctor with their statistics
type TopDoctor struct {
	DoctorID              string `json:"doctor_id"`
	DoctorName            string `json:"doctor_name"`
	TotalAppointments     int    `json:"total_appointments"`
	CompletedAppointments int    `json:"completed_appointments"`
}

// TopService represents a service with usage statistics
type TopService struct {
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	TotalCitas  int    `json:"total_citas"`
}

// CancellationStats represents cancellation statistics
type CancellationStats struct {
	TotalCancelled   int                  `json:"total_cancelled"`
	CancellationRate float64              `json:"cancellation_rate"`
	TopReasons       []CancellationReason `json:"top_reasons"`
}

// CancellationReason represents a cancellation reason with count
type CancellationReason struct {
	Reason string `json:"reason"`
	Count  int    `json:"count"`
}

// PeakHoursStats represents statistics about peak hours
type PeakHoursStats struct {
	Hour  int `json:"hour"`  // 0-23
	Count int `json:"count"`
}

// DateRangeRequest represents a date range filter
type DateRangeRequest struct {
	StartDate time.Time
	EndDate   time.Time
}

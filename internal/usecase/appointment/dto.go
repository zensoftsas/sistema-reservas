package appointment

import "time"

// CreateAppointmentRequest represents the input data for creating a new appointment
type CreateAppointmentRequest struct {
	DoctorID        string `json:"doctor_id"`
	AppointmentDate string `json:"appointment_date"` // Format: "2006-01-02" (YYYY-MM-DD)
	AppointmentTime string `json:"appointment_time"` // Format: "15:04" (HH:MM 24-hour)
	Reason          string `json:"reason"`
}

// CreateAppointmentResponse represents the output data after successfully creating an appointment
type CreateAppointmentResponse struct {
	ID              string    `json:"id"`
	PatientID       string    `json:"patient_id"`
	DoctorID        string    `json:"doctor_id"`
	AppointmentDate string    `json:"appointment_date"`
	AppointmentTime string    `json:"appointment_time"`
	Status          string    `json:"status"`
	Reason          string    `json:"reason"`
	CreatedAt       time.Time `json:"created_at"`
}

// GetAppointmentResponse represents the output data for getting an appointment by ID
type GetAppointmentResponse struct {
	ID              string    `json:"id"`
	PatientID       string    `json:"patient_id"`
	DoctorID        string    `json:"doctor_id"`
	AppointmentDate string    `json:"appointment_date"`
	AppointmentTime string    `json:"appointment_time"`
	Status          string    `json:"status"`
	Reason          string    `json:"reason"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
}

// CancelAppointmentRequest represents the input data for canceling an appointment
type CancelAppointmentRequest struct {
	Reason string `json:"reason"` // Cancellation reason
}

// ConfirmAppointmentResponse represents the response after confirming an appointment
type ConfirmAppointmentResponse struct {
	ID              string    `json:"id"`
	PatientID       string    `json:"patient_id"`
	DoctorID        string    `json:"doctor_id"`
	AppointmentDate string    `json:"appointment_date"`
	AppointmentTime string    `json:"appointment_time"`
	Status          string    `json:"status"`
	Reason          string    `json:"reason"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CompleteAppointmentRequest represents the input for completing an appointment
type CompleteAppointmentRequest struct {
	Notes string `json:"notes"` // Medical notes from the consultation
}

// CompleteAppointmentResponse represents the response after completing an appointment
type CompleteAppointmentResponse struct {
	ID              string    `json:"id"`
	PatientID       string    `json:"patient_id"`
	DoctorID        string    `json:"doctor_id"`
	AppointmentDate string    `json:"appointment_date"`
	AppointmentTime string    `json:"appointment_time"`
	Status          string    `json:"status"`
	Reason          string    `json:"reason"`
	Notes           string    `json:"notes"`
	UpdatedAt       time.Time `json:"updated_at"`
}

package analytics

import (
	"context"

	"version-1-0/internal/repository"
)

// GetTopDoctorsUseCase handles retrieving top doctors statistics
type GetTopDoctorsUseCase struct {
	appointmentRepo repository.AppointmentRepository
	userRepo        repository.UserRepository
}

// NewGetTopDoctorsUseCase creates a new instance
func NewGetTopDoctorsUseCase(
	appointmentRepo repository.AppointmentRepository,
	userRepo repository.UserRepository,
) *GetTopDoctorsUseCase {
	return &GetTopDoctorsUseCase{
		appointmentRepo: appointmentRepo,
		userRepo:        userRepo,
	}
}

// Execute retrieves top doctors by number of appointments
func (uc *GetTopDoctorsUseCase) Execute(ctx context.Context, limit int) ([]*TopDoctor, error) {
	if limit <= 0 {
		limit = 10 // Default to top 10
	}

	// Get top doctors from repository
	doctorsData, err := uc.appointmentRepo.GetTopDoctors(ctx, limit)
	if err != nil {
		return nil, err
	}

	var results []*TopDoctor
	for _, data := range doctorsData {
		// Get doctor user info to get name
		// Note: doctor_id in appointments is doctor.id, need to get user info
		// For now, we'll use doctor_id as identifier
		// In production, you'd join with doctors and users tables

		results = append(results, &TopDoctor{
			DoctorID:              data.DoctorID,
			DoctorName:            "Doctor " + data.DoctorID[:8], // Placeholder
			TotalAppointments:     data.TotalAppointments,
			CompletedAppointments: data.CompletedAppointments,
		})
	}

	return results, nil
}

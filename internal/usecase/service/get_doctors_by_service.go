package service

import (
	"context"
	"errors"

	"version-1-0/internal/repository"
	"version-1-0/internal/usecase/user"
)

// GetDoctorsByServiceUseCase handles retrieving doctors that offer a specific service
type GetDoctorsByServiceUseCase struct {
	doctorServiceRepo repository.DoctorServiceRepository
	serviceRepo       repository.ServiceRepository
}

// NewGetDoctorsByServiceUseCase creates a new instance
func NewGetDoctorsByServiceUseCase(
	doctorServiceRepo repository.DoctorServiceRepository,
	serviceRepo repository.ServiceRepository,
) *GetDoctorsByServiceUseCase {
	return &GetDoctorsByServiceUseCase{
		doctorServiceRepo: doctorServiceRepo,
		serviceRepo:       serviceRepo,
	}
}

// Execute retrieves all doctors offering a specific service
func (uc *GetDoctorsByServiceUseCase) Execute(ctx context.Context, serviceID string) ([]user.GetUserResponse, error) {
	// Validate input
	if serviceID == "" {
		return nil, errors.New("service ID is required")
	}

	// Verify service exists
	service, err := uc.serviceRepo.FindByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	if service == nil {
		return nil, errors.New("service not found")
	}

	// Get doctors
	doctors, err := uc.doctorServiceRepo.FindDoctorsByService(ctx, serviceID)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	responses := make([]user.GetUserResponse, len(doctors))
	for i, doctor := range doctors {
		responses[i] = user.GetUserResponse{
			ID:        doctor.ID,
			Email:     doctor.Email,
			FirstName: doctor.FirstName,
			LastName:  doctor.LastName,
			Phone:     doctor.Phone,
			Role:      string(doctor.Role),
			IsActive:  doctor.IsActive,
			CreatedAt: doctor.CreatedAt,
		}
	}

	return responses, nil
}

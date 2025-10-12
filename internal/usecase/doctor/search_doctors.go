package doctor

import (
	"context"
	"strings"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// SearchDoctorsUseCase handles searching for doctors by specialty
type SearchDoctorsUseCase struct {
	userRepo repository.UserRepository
}

// NewSearchDoctorsUseCase creates a new instance of SearchDoctorsUseCase
func NewSearchDoctorsUseCase(userRepo repository.UserRepository) *SearchDoctorsUseCase {
	return &SearchDoctorsUseCase{
		userRepo: userRepo,
	}
}

// Execute searches for doctors, optionally filtered by specialty
// Uses the first 6 characters of the search term to avoid accent/encoding issues
// This makes the search work with accented input (e.g., "Cardiología" finds "Cardiologia")
func (uc *SearchDoctorsUseCase) Execute(ctx context.Context, specialty string) ([]DoctorSearchResponse, error) {
	var users []*domain.User
	var err error

	// If specialty is provided, filter by specialty
	if strings.TrimSpace(specialty) != "" {
		searchTerm := strings.TrimSpace(specialty)

		// Use only first 6 characters to avoid accent issues
		// This allows "Cardiología" to match "Cardiologia"
		if len(searchTerm) > 6 {
			searchTerm = searchTerm[:6]
		}

		users, err = uc.userRepo.FindDoctorsBySpecialty(ctx, searchTerm)
	} else {
		// Otherwise, get all doctors
		users, err = uc.userRepo.GetAllDoctors(ctx)
	}

	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	responses := make([]DoctorSearchResponse, len(users))
	for i, user := range users {
		responses[i] = DoctorSearchResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		}
	}

	return responses, nil
}

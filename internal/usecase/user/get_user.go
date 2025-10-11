package user

import (
	"context"
	"errors"
	"strings"

	"version-1-0/internal/repository"
)

// GetUserUseCase handles the business logic for retrieving a user by ID
type GetUserUseCase struct {
	userRepo repository.UserRepository
}

// NewGetUserUseCase creates a new instance of GetUserUseCase
func NewGetUserUseCase(userRepo repository.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves a user by their ID
// Returns the user information or an error if the user is not found
func (uc *GetUserUseCase) Execute(ctx context.Context, id string) (*GetUserResponse, error) {
	// Validate ID is not empty
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("user ID is required")
	}

	// Find user by ID
	user, err := uc.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Build and return response
	response := &GetUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
	}

	return response, nil
}

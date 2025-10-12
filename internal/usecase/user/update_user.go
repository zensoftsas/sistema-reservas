package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"version-1-0/internal/repository"
)

// UpdateUserUseCase handles the business logic for updating user information
type UpdateUserUseCase struct {
	userRepo repository.UserRepository
}

// NewUpdateUserUseCase creates a new instance of UpdateUserUseCase
func NewUpdateUserUseCase(userRepo repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo: userRepo,
	}
}

// Execute updates a user's information with permission validation
// Only admins can update any user, regular users can only update themselves
func (uc *UpdateUserUseCase) Execute(ctx context.Context, userID string, authenticatedUserID string, authenticatedUserRole string, req UpdateUserRequest) (*UpdateUserResponse, error) {
	// Validate permissions: only admins or the user themselves can update
	if authenticatedUserRole != "admin" && userID != authenticatedUserID {
		return nil, errors.New("insufficient permissions to update this user")
	}

	// Retrieve existing user
	existingUser, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	// Validate and update first name if provided
	if strings.TrimSpace(req.FirstName) != "" {
		trimmedFirstName := strings.TrimSpace(req.FirstName)
		if len(trimmedFirstName) == 0 {
			return nil, errors.New("first name cannot be empty")
		}
		existingUser.FirstName = trimmedFirstName
	}

	// Validate and update last name if provided
	if strings.TrimSpace(req.LastName) != "" {
		trimmedLastName := strings.TrimSpace(req.LastName)
		if len(trimmedLastName) == 0 {
			return nil, errors.New("last name cannot be empty")
		}
		existingUser.LastName = trimmedLastName
	}

	// Update phone if provided
	if strings.TrimSpace(req.Phone) != "" {
		existingUser.Phone = strings.TrimSpace(req.Phone)
	}

	// Update password if provided
	if req.Password != "" {
		// Validate password length
		if len(req.Password) < 8 {
			return nil, errors.New("password must be at least 8 characters long")
		}

		// Hash the new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}

		existingUser.PasswordHash = string(hashedPassword)
	}

	// Update timestamp
	existingUser.UpdatedAt = time.Now()

	// Save to database
	err = uc.userRepo.Update(ctx, existingUser)
	if err != nil {
		return nil, err
	}

	// Build and return response
	response := &UpdateUserResponse{
		ID:        existingUser.ID,
		Email:     existingUser.Email,
		FirstName: existingUser.FirstName,
		LastName:  existingUser.LastName,
		Phone:     existingUser.Phone,
		Role:      string(existingUser.Role),
		IsActive:  existingUser.IsActive,
		UpdatedAt: existingUser.UpdatedAt,
	}

	return response, nil
}

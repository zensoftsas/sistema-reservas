package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// CreateUserUseCase handles the business logic for creating a new user
type CreateUserUseCase struct {
	userRepo repository.UserRepository
}

// NewCreateUserUseCase creates a new instance of CreateUserUseCase
func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
	}
}

// Execute creates a new user with the provided data
// Returns the created user information or an error if validation or creation fails
func (uc *CreateUserUseCase) Execute(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
	// Validate email
	if strings.TrimSpace(req.Email) == "" {
		return nil, errors.New("email is required")
	}

	// Validate password length
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	// Validate first name
	if strings.TrimSpace(req.FirstName) == "" {
		return nil, errors.New("first name is required")
	}

	// Validate last name
	if strings.TrimSpace(req.LastName) == "" {
		return nil, errors.New("last name is required")
	}

	// Validate role
	if !domain.IsValidRole(req.Role) {
		return nil, errors.New("invalid role: must be admin, doctor, or patient")
	}

	// Check if email already exists
	existingUser, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user entity
	now := time.Now()
	user := domain.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Role:         domain.UserRole(req.Role),
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Validate the user entity
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Save user to repository
	if err := uc.userRepo.Create(ctx, &user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Return response without password
	response := &CreateUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}

	return response, nil
}

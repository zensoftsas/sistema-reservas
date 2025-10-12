package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"version-1-0/internal/repository"
)

// LoginUseCase handles the business logic for user authentication
type LoginUseCase struct {
	userRepo           repository.UserRepository
	jwtSecret          string
	jwtExpirationHours int
}

// NewLoginUseCase creates a new instance of LoginUseCase
func NewLoginUseCase(userRepo repository.UserRepository, jwtSecret string, jwtExpirationHours int) *LoginUseCase {
	return &LoginUseCase{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		jwtExpirationHours: jwtExpirationHours,
	}
}

// Execute authenticates a user and returns a JWT token if credentials are valid
// Returns an error if credentials are invalid or user is inactive
func (uc *LoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Validate email is not empty
	if strings.TrimSpace(req.Email) == "" {
		return nil, errors.New("email is required")
	}

	// Validate password is not empty
	if strings.TrimSpace(req.Password) == "" {
		return nil, errors.New("password is required")
	}

	// Find user by email
	user, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		// Return generic error for security (don't reveal if email exists)
		return nil, errors.New("invalid credentials")
	}

	// Verify password with bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		// Password doesn't match
		return nil, errors.New("invalid credentials")
	}

	// Verify user is active
	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	// Generate JWT token
	expiresAt := time.Now().Add(time.Duration(uc.jwtExpirationHours) * time.Hour)

	// Create JWT claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    string(user.Role),
		"exp":     expiresAt.Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	tokenString, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Build and return login response
	response := &LoginResponse{
		Token:     tokenString,
		ExpiresAt: expiresAt,
	}

	// Set user information
	response.User.ID = user.ID
	response.User.Email = user.Email
	response.User.FirstName = user.FirstName
	response.User.LastName = user.LastName
	response.User.Role = string(user.Role)

	return response, nil
}

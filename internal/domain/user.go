package domain

import (
	"errors"
	"strings"
	"time"
)

// UserRole represents the role of a user in the system
type UserRole string

// User role constants
const (
	RoleAdmin   UserRole = "admin"
	RoleDoctor  UserRole = "doctor"
	RolePatient UserRole = "patient"
)

// User represents a user entity in the medical reservation system
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password hash in JSON
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Phone        string    `json:"phone"`
	Role         UserRole  `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Validate checks if the User entity has all required fields properly set
// Returns an error if any validation rule fails
func (u *User) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return errors.New("user ID is required")
	}

	if strings.TrimSpace(u.Email) == "" {
		return errors.New("user email is required")
	}

	if strings.TrimSpace(u.PasswordHash) == "" {
		return errors.New("user password hash is required")
	}

	if strings.TrimSpace(u.FirstName) == "" {
		return errors.New("user first name is required")
	}

	if strings.TrimSpace(u.LastName) == "" {
		return errors.New("user last name is required")
	}

	if u.Role != RoleAdmin && u.Role != RoleDoctor && u.Role != RolePatient {
		return errors.New("invalid user role")
	}

	if u.CreatedAt.IsZero() {
		return errors.New("user created at is required")
	}

	if u.UpdatedAt.IsZero() {
		return errors.New("user updated at is required")
	}

	return nil
}

// FullName returns the user's full name combining first and last name
func (u *User) FullName() string {
	return strings.TrimSpace(u.FirstName + " " + u.LastName)
}

// IsValidRole checks if a given role string is a valid UserRole
func IsValidRole(role string) bool {
	r := UserRole(role)
	return r == RoleAdmin || r == RoleDoctor || r == RolePatient
}

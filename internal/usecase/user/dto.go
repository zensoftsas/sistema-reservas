package user

import "time"

// CreateUserRequest represents the input data for creating a new user
type CreateUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Role      string `json:"role" validate:"required"`
}

// CreateUserResponse represents the output data after successfully creating a user
// Password is intentionally excluded for security reasons
type CreateUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// GetUserResponse represents the output data for getting a user by ID
type GetUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// ListUsersRequest represents the input data for listing users with pagination
type ListUsersRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ListUsersResponse represents the output data for listing users
type ListUsersResponse struct {
	Users   []GetUserResponse `json:"users"`
	Total   int               `json:"total"`
	Limit   int               `json:"limit"`
	Offset  int               `json:"offset"`
	HasMore bool              `json:"has_more"`
}

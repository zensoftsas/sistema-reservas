package user

import (
	"context"

	"version-1-0/internal/repository"
)

// ListUsersUseCase handles the business logic for listing users with pagination
type ListUsersUseCase struct {
	userRepo repository.UserRepository
}

// NewListUsersUseCase creates a new instance of ListUsersUseCase
func NewListUsersUseCase(userRepo repository.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves a paginated list of users
// Validates pagination parameters and applies reasonable defaults and limits
func (uc *ListUsersUseCase) Execute(ctx context.Context, req ListUsersRequest) (*ListUsersResponse, error) {
	// Validate and set default limit
	if req.Limit <= 0 {
		req.Limit = 20 // default limit
	}

	// Validate offset
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Enforce maximum limit
	if req.Limit > 100 {
		req.Limit = 100 // maximum limit
	}

	// Retrieve users from repository
	users, err := uc.userRepo.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	// Convert domain users to response DTOs
	userResponses := make([]GetUserResponse, len(users))
	for i, user := range users {
		userResponses[i] = GetUserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		}
	}

	// Calculate if there are more users available
	// If we got exactly the limit requested, there might be more
	hasMore := len(users) == req.Limit

	// Build response
	// Note: Total is currently set to len(users) for simplicity
	// In production, this should be a separate COUNT query to get the actual total
	response := &ListUsersResponse{
		Users:   userResponses,
		Total:   len(users), // TODO: Replace with actual count from database
		Limit:   req.Limit,
		Offset:  req.Offset,
		HasMore: hasMore,
	}

	return response, nil
}

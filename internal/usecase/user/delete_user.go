package user

import (
	"context"
	"errors"

	"version-1-0/internal/repository"
)

// DeleteUserUseCase handles the business logic for deleting (deactivating) users
type DeleteUserUseCase struct {
	userRepo repository.UserRepository
}

// NewDeleteUserUseCase creates a new instance of DeleteUserUseCase
func NewDeleteUserUseCase(userRepo repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepo: userRepo,
	}
}

// Execute performs a soft delete on a user (marks as inactive)
// Only administrators can delete users
// Admins cannot delete their own account to prevent lockout
func (uc *DeleteUserUseCase) Execute(ctx context.Context, userID string, authenticatedUserID string, authenticatedUserRole string) error {
	// Validate that only admins can delete users
	if authenticatedUserRole != "admin" {
		return errors.New("only administrators can delete users")
	}

	// Prevent admins from deleting their own account (prevents lockout)
	if userID == authenticatedUserID {
		return errors.New("cannot delete your own account")
	}

	// Verify that the user exists
	existingUser, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	// Verify that the user is currently active
	if !existingUser.IsActive {
		return errors.New("user is already inactive")
	}

	// Perform soft delete (marks user as inactive)
	err = uc.userRepo.Delete(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

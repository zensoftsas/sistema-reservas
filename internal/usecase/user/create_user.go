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
	"version-1-0/internal/repository/sqlite"
)

// CreateUserUseCase handles the business logic for creating a new user
type CreateUserUseCase struct {
	userRepo    repository.UserRepository
	doctorRepo  repository.DoctorRepository
	patientRepo repository.PatientRepository
}

// NewCreateUserUseCase creates a new instance of CreateUserUseCase
func NewCreateUserUseCase(
	userRepo repository.UserRepository,
	doctorRepo repository.DoctorRepository,
	patientRepo repository.PatientRepository,
) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:    userRepo,
		doctorRepo:  doctorRepo,
		patientRepo: patientRepo,
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
	userID := uuid.New().String()
	user := domain.User{
		ID:           userID,
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

	// If user is admin, use simple creation without transaction
	if user.Role == domain.RoleAdmin {
		if err := uc.userRepo.Create(ctx, &user); err != nil {
			return nil, errors.New("failed to create user")
		}
	} else if user.Role == domain.RoleDoctor {
		// For doctors, use transaction to ensure atomicity
		// Get DB connection from repository (cast to concrete type)
		userRepoImpl, ok := uc.userRepo.(*sqlite.SqliteUserRepository)
		if !ok {
			return nil, errors.New("failed to get database connection")
		}

		db := userRepoImpl.GetDB()

		// Begin transaction
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return nil, errors.New("failed to begin transaction")
		}

		// Ensure rollback on error
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		// Create user within transaction
		if err = userRepoImpl.CreateWithTx(ctx, tx, &user); err != nil {
			return nil, errors.New("failed to create user")
		}

		// Generate unique license number using last 6 characters of user ID
		licenseNumber := "LIC-" + userID[len(userID)-6:]

		doctor := domain.Doctor{
			ID:                 uuid.New().String(),
			UserID:             userID,
			Specialty:          "Medicina General", // Default specialty
			LicenseNumber:      licenseNumber,
			YearsOfExperience:  0,
			Education:          "",
			Bio:                "",
			ConsultationFee:    0.0,
			IsAvailable:        true,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		// Validate doctor profile
		if err = doctor.Validate(); err != nil {
			return nil, errors.New("failed to validate doctor profile: " + err.Error())
		}

		// Create doctor profile within same transaction
		doctorRepoImpl, ok := uc.doctorRepo.(*sqlite.SqliteDoctorRepository)
		if !ok {
			return nil, errors.New("failed to get doctor repository")
		}

		if err = doctorRepoImpl.CreateWithTx(ctx, tx, &doctor); err != nil {
			return nil, errors.New("failed to create doctor profile")
		}

		// Commit transaction - if this succeeds, both user and doctor are created atomically
		if err = tx.Commit(); err != nil {
			return nil, errors.New("failed to commit transaction")
		}
	} else if user.Role == domain.RolePatient {
		// For patients, use transaction to ensure atomicity
		userRepoImpl, ok := uc.userRepo.(*sqlite.SqliteUserRepository)
		if !ok {
			return nil, errors.New("failed to get database connection")
		}

		db := userRepoImpl.GetDB()

		// Begin transaction
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return nil, errors.New("failed to begin transaction")
		}

		// Ensure rollback on error
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		// Create user within transaction
		if err = userRepoImpl.CreateWithTx(ctx, tx, &user); err != nil {
			return nil, errors.New("failed to create user")
		}

		// Create default patient profile
		// Use a reasonable default birthdate (18 years ago)
		defaultBirthdate := time.Now().AddDate(-18, 0, 0)

		patient := domain.Patient{
			ID:                    uuid.New().String(),
			UserID:                userID,
			Birthdate:             defaultBirthdate,
			DocumentType:          "DNI",        // Default document type
			DocumentNumber:        "00000000",   // Placeholder - should be updated later
			Address:               "Por definir", // Placeholder
			EmergencyContactName:  "Por definir", // Placeholder
			EmergencyContactPhone: req.Phone,    // Use patient's phone as default
			BloodType:             "",           // Optional
			Allergies:             []string{},   // Empty by default
			CreatedAt:             now,
			UpdatedAt:             now,
		}

		// Validate patient profile
		if err = patient.Validate(); err != nil {
			return nil, errors.New("failed to validate patient profile: " + err.Error())
		}

		// Create patient profile within same transaction
		patientRepoImpl, ok := uc.patientRepo.(*sqlite.SqlitePatientRepository)
		if !ok {
			return nil, errors.New("failed to get patient repository")
		}

		if err = patientRepoImpl.CreateWithTx(ctx, tx, &patient); err != nil {
			return nil, errors.New("failed to create patient profile")
		}

		// Commit transaction - if this succeeds, both user and patient are created atomically
		if err = tx.Commit(); err != nil {
			return nil, errors.New("failed to commit transaction")
		}
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

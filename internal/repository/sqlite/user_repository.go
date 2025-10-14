package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// SqliteUserRepository implements the UserRepository interface using SQLite
type SqliteUserRepository struct {
	db *sql.DB
}

// NewSqliteUserRepository creates a new instance of SqliteUserRepository
func NewSqliteUserRepository(db *sql.DB) repository.UserRepository {
	return &SqliteUserRepository{
		db: db,
	}
}

// Create inserts a new user into the database
func (r *SqliteUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, phone, role, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
		user.IsActive,
		user.CreatedAt.Format(time.RFC3339),
		user.UpdatedAt.Format(time.RFC3339),
	)

	if err != nil {
		return err
	}

	return nil
}

// FindByID retrieves a user by their unique identifier
// Returns nil if the user is not found
func (r *SqliteUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, role, is_active, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user domain.User
	var isActive int
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&isActive,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Convert is_active from int to bool
	user.IsActive = isActive == 1

	// Parse timestamps
	user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, err
	}

	user.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail retrieves a user by their email address
// Returns nil if the user is not found
func (r *SqliteUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, role, is_active, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var user domain.User
	var isActive int
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Role,
		&isActive,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Convert is_active from int to bool
	user.IsActive = isActive == 1

	// Parse timestamps
	user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, err
	}

	user.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update modifies an existing user in the database
func (r *SqliteUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = ?, password_hash = ?, first_name = ?, last_name = ?, phone = ?, role = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Role,
		user.IsActive,
		user.UpdatedAt.Format(time.RFC3339),
		user.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Delete performs a soft delete by marking the user as inactive
func (r *SqliteUserRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET is_active = 0,
		    updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(ctx, query, time.Now().Format(time.RFC3339), id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// List retrieves a paginated list of users
func (r *SqliteUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, role, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var user domain.User
		var isActive int
		var createdAt, updatedAt string

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.Role,
			&isActive,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Convert is_active from int to bool
		user.IsActive = isActive == 1

		// Parse timestamps
		user.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, err
		}

		user.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// FindDoctorsBySpecialty retrieves all active doctors filtered by specialty
func (r *SqliteUserRepository) FindDoctorsBySpecialty(ctx context.Context, specialty string) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.phone, u.role, u.is_active, u.created_at, u.updated_at
		FROM users u
		INNER JOIN doctors d ON u.id = d.user_id
		WHERE u.role = 'doctor'
		AND u.is_active = 1
		AND LOWER(d.specialty) LIKE LOWER(?)
		ORDER BY u.last_name ASC, u.first_name ASC
	`

	// Use LIKE with % for flexible matching
	searchPattern := "%" + specialty + "%"

	rows, err := r.db.QueryContext(ctx, query, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var createdAt, updatedAt string

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.Role,
			&user.IsActive,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		users = append(users, &user)
	}

	return users, rows.Err()
}

// GetAllDoctors retrieves all active doctors that have a complete profile
func (r *SqliteUserRepository) GetAllDoctors(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.email, u.password_hash, u.first_name, u.last_name, u.phone, u.role, u.is_active, u.created_at, u.updated_at
		FROM users u
		INNER JOIN doctors d ON u.id = d.user_id
		WHERE u.role = 'doctor' AND u.is_active = 1
		ORDER BY u.last_name ASC, u.first_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		var createdAt, updatedAt string

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&user.Role,
			&user.IsActive,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		user.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		users = append(users, &user)
	}

	return users, rows.Err()
}

// FindDoctorIDByUserID returns the doctor.id for a given user_id
func (r *SqliteUserRepository) FindDoctorIDByUserID(ctx context.Context, userID string) (string, error) {
	query := `SELECT id FROM doctors WHERE user_id = ?`

	var doctorID string
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&doctorID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("doctor not found")
		}
		return "", err
	}

	return doctorID, nil
}

// CountByRole counts users by role
func (r *SqliteUserRepository) CountByRole(ctx context.Context, role string) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE role = ? AND is_active = 1`

	var count int
	err := r.db.QueryRowContext(ctx, query, role).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountAllActive counts all active users
func (r *SqliteUserRepository) CountAllActive(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE is_active = 1`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

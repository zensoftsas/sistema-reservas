package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// SqliteServiceRepository implements the ServiceRepository interface using SQLite
type SqliteServiceRepository struct {
	db *sql.DB
}

// NewSqliteServiceRepository creates a new instance of SqliteServiceRepository
func NewSqliteServiceRepository(db *sql.DB) repository.ServiceRepository {
	return &SqliteServiceRepository{
		db: db,
	}
}

// Create inserts a new service into the database
func (r *SqliteServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	query := `
		INSERT INTO services (id, name, description, duration_minutes, price, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		service.ID,
		service.Name,
		service.Description,
		service.DurationMinutes,
		service.Price,
		service.IsActive,
		service.CreatedAt,
		service.UpdatedAt,
	)

	return err
}

// FindByID retrieves a service by its unique identifier
func (r *SqliteServiceRepository) FindByID(ctx context.Context, id string) (*domain.Service, error) {
	query := `
		SELECT id, name, description, duration_minutes, price, is_active, created_at, updated_at
		FROM services
		WHERE id = $1
	`

	var service domain.Service
	var isActive bool
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		&service.Description,
		&service.DurationMinutes,
		&service.Price,
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

	// Assign scanned values
	service.IsActive = isActive
	service.CreatedAt = createdAt
	service.UpdatedAt = updatedAt

	return &service, nil
}

// ListActive retrieves all active services
func (r *SqliteServiceRepository) ListActive(ctx context.Context) ([]*domain.Service, error) {
	query := `
		SELECT id, name, description, duration_minutes, price, is_active, created_at, updated_at
		FROM services
		WHERE is_active = TRUE
		ORDER BY name ASC
	`

	return r.queryServices(ctx, query)
}

// ListAll retrieves all services (active and inactive)
func (r *SqliteServiceRepository) ListAll(ctx context.Context) ([]*domain.Service, error) {
	query := `
		SELECT id, name, description, duration_minutes, price, is_active, created_at, updated_at
		FROM services
		ORDER BY name ASC
	`

	return r.queryServices(ctx, query)
}

// Update modifies an existing service in the database
func (r *SqliteServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	query := `
		UPDATE services
		SET name = $1, description = $2, duration_minutes = $3, price = $4, is_active = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		service.Name,
		service.Description,
		service.DurationMinutes,
		service.Price,
		service.IsActive,
		service.UpdatedAt,
		service.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("service not found")
	}

	return nil
}

// Delete removes a service from the database
func (r *SqliteServiceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM services WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("service not found")
	}

	return nil
}

// queryServices is a helper method to query services
func (r *SqliteServiceRepository) queryServices(ctx context.Context, query string, args ...interface{}) ([]*domain.Service, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []*domain.Service

	for rows.Next() {
		var service domain.Service
		var isActive bool
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Description,
			&service.DurationMinutes,
			&service.Price,
			&isActive,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Assign scanned values
		service.IsActive = isActive
		service.CreatedAt = createdAt
		service.UpdatedAt = updatedAt

		services = append(services, &service)
	}

	return services, rows.Err()
}

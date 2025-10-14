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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
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
		service.CreatedAt.Format(time.RFC3339),
		service.UpdatedAt.Format(time.RFC3339),
	)

	return err
}

// FindByID retrieves a service by its unique identifier
func (r *SqliteServiceRepository) FindByID(ctx context.Context, id string) (*domain.Service, error) {
	query := `
		SELECT id, name, description, duration_minutes, price, is_active, created_at, updated_at
		FROM services
		WHERE id = ?
	`

	var service domain.Service
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.Name,
		&service.Description,
		&service.DurationMinutes,
		&service.Price,
		&service.IsActive,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Parse timestamps
	service.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	service.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &service, nil
}

// ListActive retrieves all active services
func (r *SqliteServiceRepository) ListActive(ctx context.Context) ([]*domain.Service, error) {
	query := `
		SELECT id, name, description, duration_minutes, price, is_active, created_at, updated_at
		FROM services
		WHERE is_active = 1
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
		SET name = ?, description = ?, duration_minutes = ?, price = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		service.Name,
		service.Description,
		service.DurationMinutes,
		service.Price,
		service.IsActive,
		service.UpdatedAt.Format(time.RFC3339),
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
	query := `DELETE FROM services WHERE id = ?`

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
		var createdAt, updatedAt string

		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Description,
			&service.DurationMinutes,
			&service.Price,
			&service.IsActive,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Parse timestamps
		service.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		service.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		services = append(services, &service)
	}

	return services, rows.Err()
}

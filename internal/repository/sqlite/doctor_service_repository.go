package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// SqliteDoctorServiceRepository implements the DoctorServiceRepository interface using SQLite
type SqliteDoctorServiceRepository struct {
	db *sql.DB
}

// NewSqliteDoctorServiceRepository creates a new instance of SqliteDoctorServiceRepository
func NewSqliteDoctorServiceRepository(db *sql.DB) repository.DoctorServiceRepository {
	return &SqliteDoctorServiceRepository{
		db: db,
	}
}

// Assign creates a relationship between a doctor and a service
func (r *SqliteDoctorServiceRepository) Assign(ctx context.Context, doctorService *domain.DoctorService) error {
	query := `
		INSERT INTO doctor_services (id, doctor_id, service_id, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		doctorService.ID,
		doctorService.DoctorID,
		doctorService.ServiceID,
		doctorService.IsActive,
		doctorService.CreatedAt,
		doctorService.UpdatedAt,
	)

	return err
}

// Remove removes a service assignment from a doctor
func (r *SqliteDoctorServiceRepository) Remove(ctx context.Context, doctorID, serviceID string) error {
	query := `DELETE FROM doctor_services WHERE doctor_id = $1 AND service_id = $2`

	result, err := r.db.ExecContext(ctx, query, doctorID, serviceID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("doctor-service relationship not found")
	}

	return nil
}

// FindDoctorsByService returns all doctors that offer a specific service
func (r *SqliteDoctorServiceRepository) FindDoctorsByService(ctx context.Context, serviceID string) ([]*domain.User, error) {
	query := `
		SELECT
			u.id,
			u.email,
			u.password_hash,
			u.first_name,
			u.last_name,
			u.phone,
			u.role,
			u.is_active,
			u.created_at,
			u.updated_at
		FROM users u
		INNER JOIN doctors d ON d.user_id = u.id
		INNER JOIN doctor_services ds ON ds.doctor_id = d.id
		WHERE ds.service_id = $1 AND ds.is_active = TRUE AND u.is_active = TRUE
		ORDER BY u.first_name ASC, u.last_name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var user domain.User
		var isActive bool
		var createdAt, updatedAt time.Time

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

		// Assign scanned values
		user.IsActive = isActive
		user.CreatedAt = createdAt
		user.UpdatedAt = updatedAt

		users = append(users, &user)
	}

	return users, rows.Err()
}

// FindServicesByDoctor returns all services offered by a specific doctor
func (r *SqliteDoctorServiceRepository) FindServicesByDoctor(ctx context.Context, doctorID string) ([]*domain.Service, error) {
	query := `
		SELECT
			s.id,
			s.name,
			s.description,
			s.duration_minutes,
			s.price,
			s.is_active,
			s.created_at,
			s.updated_at
		FROM services s
		INNER JOIN doctor_services ds ON ds.service_id = s.id
		WHERE ds.doctor_id = $1 AND ds.is_active = TRUE AND s.is_active = TRUE
		ORDER BY s.name ASC
	`

	rows, err := r.db.QueryContext(ctx, query, doctorID)
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

// IsAssigned checks if a doctor is assigned to a service
func (r *SqliteDoctorServiceRepository) IsAssigned(ctx context.Context, doctorID, serviceID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM doctor_services
		WHERE doctor_id = $1 AND service_id = $2 AND is_active = TRUE
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, doctorID, serviceID).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// FindByDoctorAndService retrieves a specific doctor-service relationship
func (r *SqliteDoctorServiceRepository) FindByDoctorAndService(ctx context.Context, doctorID, serviceID string) (*domain.DoctorService, error) {
	query := `
		SELECT id, doctor_id, service_id, is_active, created_at, updated_at
		FROM doctor_services
		WHERE doctor_id = $1 AND service_id = $2
	`

	var ds domain.DoctorService
	var isActive bool
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, doctorID, serviceID).Scan(
		&ds.ID,
		&ds.DoctorID,
		&ds.ServiceID,
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
	ds.IsActive = isActive
	ds.CreatedAt = createdAt
	ds.UpdatedAt = updatedAt

	return &ds, nil
}

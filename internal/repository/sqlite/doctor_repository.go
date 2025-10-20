package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// SqliteDoctorRepository implements the DoctorRepository interface using PostgreSQL
type SqliteDoctorRepository struct {
	db *sql.DB
}

// NewSqliteDoctorRepository creates a new instance of SqliteDoctorRepository
func NewSqliteDoctorRepository(db *sql.DB) repository.DoctorRepository {
	return &SqliteDoctorRepository{
		db: db,
	}
}

// Create inserts a new doctor into the database
func (r *SqliteDoctorRepository) Create(ctx context.Context, doctor *domain.Doctor) error {
	return r.CreateWithTx(ctx, nil, doctor)
}

// CreateWithTx inserts a new doctor using a transaction or database connection
func (r *SqliteDoctorRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, doctor *domain.Doctor) error {
	query := `
		INSERT INTO doctors (id, user_id, specialty, license_number, years_of_experience,
		                     education, bio, consultation_fee, is_available, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	var err error
	if tx != nil {
		_, err = tx.ExecContext(
			ctx,
			query,
			doctor.ID,
			doctor.UserID,
			doctor.Specialty,
			doctor.LicenseNumber,
			doctor.YearsOfExperience,
			doctor.Education,
			doctor.Bio,
			doctor.ConsultationFee,
			doctor.IsAvailable,
			doctor.CreatedAt,
			doctor.UpdatedAt,
		)
	} else {
		_, err = r.db.ExecContext(
			ctx,
			query,
			doctor.ID,
			doctor.UserID,
			doctor.Specialty,
			doctor.LicenseNumber,
			doctor.YearsOfExperience,
			doctor.Education,
			doctor.Bio,
			doctor.ConsultationFee,
			doctor.IsAvailable,
			doctor.CreatedAt,
			doctor.UpdatedAt,
		)
	}

	return err
}

// FindByID retrieves a doctor by their unique identifier
func (r *SqliteDoctorRepository) FindByID(ctx context.Context, id string) (*domain.Doctor, error) {
	query := `
		SELECT id, user_id, specialty, license_number, years_of_experience,
		       education, bio, consultation_fee, is_available, created_at, updated_at
		FROM doctors
		WHERE id = $1
	`

	var doctor domain.Doctor
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&doctor.ID,
		&doctor.UserID,
		&doctor.Specialty,
		&doctor.LicenseNumber,
		&doctor.YearsOfExperience,
		&doctor.Education,
		&doctor.Bio,
		&doctor.ConsultationFee,
		&doctor.IsAvailable,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &doctor, nil
}

// FindByUserID retrieves a doctor by their associated user ID
func (r *SqliteDoctorRepository) FindByUserID(ctx context.Context, userID string) (*domain.Doctor, error) {
	query := `
		SELECT id, user_id, specialty, license_number, years_of_experience,
		       education, bio, consultation_fee, is_available, created_at, updated_at
		FROM doctors
		WHERE user_id = $1
	`

	var doctor domain.Doctor
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&doctor.ID,
		&doctor.UserID,
		&doctor.Specialty,
		&doctor.LicenseNumber,
		&doctor.YearsOfExperience,
		&doctor.Education,
		&doctor.Bio,
		&doctor.ConsultationFee,
		&doctor.IsAvailable,
		&doctor.CreatedAt,
		&doctor.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &doctor, nil
}

// FindBySpecialty retrieves all doctors with a specific specialty
func (r *SqliteDoctorRepository) FindBySpecialty(ctx context.Context, specialty string) ([]*domain.Doctor, error) {
	query := `
		SELECT id, user_id, specialty, license_number, years_of_experience,
		       education, bio, consultation_fee, is_available, created_at, updated_at
		FROM doctors
		WHERE LOWER(specialty) LIKE LOWER($1)
		ORDER BY created_at DESC
	`

	searchPattern := "%" + specialty + "%"
	rows, err := r.db.QueryContext(ctx, query, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []*domain.Doctor
	for rows.Next() {
		var doctor domain.Doctor
		err := rows.Scan(
			&doctor.ID,
			&doctor.UserID,
			&doctor.Specialty,
			&doctor.LicenseNumber,
			&doctor.YearsOfExperience,
			&doctor.Education,
			&doctor.Bio,
			&doctor.ConsultationFee,
			&doctor.IsAvailable,
			&doctor.CreatedAt,
			&doctor.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, &doctor)
	}

	return doctors, rows.Err()
}

// Update modifies an existing doctor in the database
func (r *SqliteDoctorRepository) Update(ctx context.Context, doctor *domain.Doctor) error {
	query := `
		UPDATE doctors
		SET specialty = $1,
		    license_number = $2,
		    years_of_experience = $3,
		    education = $4,
		    bio = $5,
		    consultation_fee = $6,
		    is_available = $7,
		    updated_at = $8
		WHERE id = $9
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		doctor.Specialty,
		doctor.LicenseNumber,
		doctor.YearsOfExperience,
		doctor.Education,
		doctor.Bio,
		doctor.ConsultationFee,
		doctor.IsAvailable,
		doctor.UpdatedAt,
		doctor.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("doctor not found")
	}

	return nil
}

// Delete removes a doctor from the repository by their ID
func (r *SqliteDoctorRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM doctors WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("doctor not found")
	}

	return nil
}

// List retrieves a paginated list of doctors
func (r *SqliteDoctorRepository) List(ctx context.Context, limit, offset int) ([]*domain.Doctor, error) {
	query := `
		SELECT id, user_id, specialty, license_number, years_of_experience,
		       education, bio, consultation_fee, is_available, created_at, updated_at
		FROM doctors
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []*domain.Doctor
	for rows.Next() {
		var doctor domain.Doctor
		err := rows.Scan(
			&doctor.ID,
			&doctor.UserID,
			&doctor.Specialty,
			&doctor.LicenseNumber,
			&doctor.YearsOfExperience,
			&doctor.Education,
			&doctor.Bio,
			&doctor.ConsultationFee,
			&doctor.IsAvailable,
			&doctor.CreatedAt,
			&doctor.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, &doctor)
	}

	return doctors, rows.Err()
}

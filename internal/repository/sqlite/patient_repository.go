package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// SqlitePatientRepository implements the PatientRepository interface using PostgreSQL
type SqlitePatientRepository struct {
	db *sql.DB
}

// NewSqlitePatientRepository creates a new instance of SqlitePatientRepository
func NewSqlitePatientRepository(db *sql.DB) repository.PatientRepository {
	return &SqlitePatientRepository{
		db: db,
	}
}

// Create inserts a new patient into the database
func (r *SqlitePatientRepository) Create(ctx context.Context, patient *domain.Patient) error {
	return r.CreateWithTx(ctx, nil, patient)
}

// CreateWithTx inserts a new patient using a transaction or database connection
func (r *SqlitePatientRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, patient *domain.Patient) error {
	query := `
		INSERT INTO patients (id, user_id, birthdate, document_type, document_number,
		                      address, emergency_contact_name, emergency_contact_phone,
		                      blood_type, allergies, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	// Convert allergies slice to comma-separated string
	allergiesStr := strings.Join(patient.Allergies, ",")

	var err error
	if tx != nil {
		_, err = tx.ExecContext(
			ctx,
			query,
			patient.ID,
			patient.UserID,
			patient.Birthdate,
			patient.DocumentType,
			patient.DocumentNumber,
			patient.Address,
			patient.EmergencyContactName,
			patient.EmergencyContactPhone,
			patient.BloodType,
			allergiesStr,
			patient.CreatedAt,
			patient.UpdatedAt,
		)
	} else {
		_, err = r.db.ExecContext(
			ctx,
			query,
			patient.ID,
			patient.UserID,
			patient.Birthdate,
			patient.DocumentType,
			patient.DocumentNumber,
			patient.Address,
			patient.EmergencyContactName,
			patient.EmergencyContactPhone,
			patient.BloodType,
			allergiesStr,
			patient.CreatedAt,
			patient.UpdatedAt,
		)
	}

	return err
}

// FindByID retrieves a patient by their unique identifier
func (r *SqlitePatientRepository) FindByID(ctx context.Context, id string) (*domain.Patient, error) {
	query := `
		SELECT id, user_id, birthdate, document_type, document_number,
		       address, emergency_contact_name, emergency_contact_phone,
		       blood_type, allergies, created_at, updated_at
		FROM patients
		WHERE id = $1
	`

	var patient domain.Patient
	var allergiesStr string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&patient.ID,
		&patient.UserID,
		&patient.Birthdate,
		&patient.DocumentType,
		&patient.DocumentNumber,
		&patient.Address,
		&patient.EmergencyContactName,
		&patient.EmergencyContactPhone,
		&patient.BloodType,
		&allergiesStr,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Convert allergies string back to slice
	if allergiesStr != "" {
		patient.Allergies = strings.Split(allergiesStr, ",")
	} else {
		patient.Allergies = []string{}
	}

	return &patient, nil
}

// FindByUserID retrieves a patient by their associated user ID
func (r *SqlitePatientRepository) FindByUserID(ctx context.Context, userID string) (*domain.Patient, error) {
	query := `
		SELECT id, user_id, birthdate, document_type, document_number,
		       address, emergency_contact_name, emergency_contact_phone,
		       blood_type, allergies, created_at, updated_at
		FROM patients
		WHERE user_id = $1
	`

	var patient domain.Patient
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&patient.ID,
		&patient.UserID,
		&patient.Birthdate,
		&patient.DocumentType,
		&patient.DocumentNumber,
		&patient.Address,
		&patient.EmergencyContactName,
		&patient.EmergencyContactPhone,
		&patient.BloodType,
		&patient.Allergies,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &patient, nil
}

// Update modifies an existing patient in the database
func (r *SqlitePatientRepository) Update(ctx context.Context, patient *domain.Patient) error {
	query := `
		UPDATE patients
		SET birthdate = $1,
		    document_type = $2,
		    document_number = $3,
		    address = $4,
		    emergency_contact_name = $5,
		    emergency_contact_phone = $6,
		    blood_type = $7,
		    allergies = $8,
		    updated_at = $9
		WHERE id = $10
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		patient.Birthdate,
		patient.DocumentType,
		patient.DocumentNumber,
		patient.Address,
		patient.EmergencyContactName,
		patient.EmergencyContactPhone,
		patient.BloodType,
		patient.Allergies,
		patient.UpdatedAt,
		patient.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("patient not found")
	}

	return nil
}

// Delete removes a patient from the repository by their ID
func (r *SqlitePatientRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM patients WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("patient not found")
	}

	return nil
}

// List retrieves a paginated list of patients
func (r *SqlitePatientRepository) List(ctx context.Context, limit, offset int) ([]*domain.Patient, error) {
	query := `
		SELECT id, user_id, birthdate, document_type, document_number,
		       address, emergency_contact_name, emergency_contact_phone,
		       blood_type, allergies, created_at, updated_at
		FROM patients
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*domain.Patient
	for rows.Next() {
		var patient domain.Patient
		err := rows.Scan(
			&patient.ID,
			&patient.UserID,
			&patient.Birthdate,
			&patient.DocumentType,
			&patient.DocumentNumber,
			&patient.Address,
			&patient.EmergencyContactName,
			&patient.EmergencyContactPhone,
			&patient.BloodType,
			&patient.Allergies,
			&patient.CreatedAt,
			&patient.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		patients = append(patients, &patient)
	}

	return patients, rows.Err()
}

package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"version-1-0/internal/domain"
	"version-1-0/internal/repository"
)

// SqliteAppointmentRepository implements the AppointmentRepository interface using SQLite
type SqliteAppointmentRepository struct {
	db *sql.DB
}

// NewSqliteAppointmentRepository creates a new instance of SqliteAppointmentRepository
func NewSqliteAppointmentRepository(db *sql.DB) repository.AppointmentRepository {
	return &SqliteAppointmentRepository{
		db: db,
	}
}

// Create inserts a new appointment into the database
func (r *SqliteAppointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
	query := `
		INSERT INTO appointments (id, patient_id, doctor_id, scheduled_at, duration, status, reason, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		appointment.ID,
		appointment.PatientID,
		appointment.DoctorID,
		appointment.ScheduledAt.Format(time.RFC3339),
		appointment.Duration,
		appointment.Status,
		appointment.Reason,
		appointment.Notes,
		appointment.CreatedAt.Format(time.RFC3339),
		appointment.UpdatedAt.Format(time.RFC3339),
	)

	return err
}

// FindByID retrieves an appointment by its unique identifier
func (r *SqliteAppointmentRepository) FindByID(ctx context.Context, id string) (*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, duration, status, reason, notes, created_at, updated_at
		FROM appointments
		WHERE id = ?
	`

	var appointment domain.Appointment
	var scheduledAt, createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&appointment.ID,
		&appointment.PatientID,
		&appointment.DoctorID,
		&scheduledAt,
		&appointment.Duration,
		&appointment.Status,
		&appointment.Reason,
		&appointment.Notes,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Parse scheduled_at
	appointment.ScheduledAt, _ = time.Parse(time.RFC3339, scheduledAt)
	appointment.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	appointment.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &appointment, nil
}

// FindByPatientID retrieves all appointments for a specific patient
func (r *SqliteAppointmentRepository) FindByPatientID(ctx context.Context, patientID string) ([]*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, duration, status, reason, notes, created_at, updated_at
		FROM appointments
		WHERE patient_id = ?
		ORDER BY scheduled_at DESC
	`

	return r.queryAppointments(ctx, query, patientID)
}

// FindByDoctorID retrieves all appointments for a specific doctor
func (r *SqliteAppointmentRepository) FindByDoctorID(ctx context.Context, doctorID string) ([]*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, duration, status, reason, notes, created_at, updated_at
		FROM appointments
		WHERE doctor_id = ?
		ORDER BY scheduled_at DESC
	`

	return r.queryAppointments(ctx, query, doctorID)
}

// FindByDoctorAndDate retrieves all appointments for a doctor on a specific date
func (r *SqliteAppointmentRepository) FindByDoctorAndDate(ctx context.Context, doctorID string, date time.Time) ([]*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, duration, status, reason, notes, created_at, updated_at
		FROM appointments
		WHERE doctor_id = ? AND DATE(scheduled_at) = DATE(?)
		ORDER BY scheduled_at ASC
	`

	return r.queryAppointments(ctx, query, doctorID, date.Format(time.RFC3339))
}

// Update modifies an existing appointment in the database
func (r *SqliteAppointmentRepository) Update(ctx context.Context, appointment *domain.Appointment) error {
	query := `
		UPDATE appointments
		SET status = ?, notes = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		appointment.Status,
		appointment.Notes,
		appointment.UpdatedAt.Format(time.RFC3339),
		appointment.ID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("appointment not found")
	}

	return nil
}

// Delete removes an appointment from the database
func (r *SqliteAppointmentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM appointments WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("appointment not found")
	}

	return nil
}

// queryAppointments is a helper method to query appointments
func (r *SqliteAppointmentRepository) queryAppointments(ctx context.Context, query string, args ...interface{}) ([]*domain.Appointment, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment

	for rows.Next() {
		var appointment domain.Appointment
		var scheduledAt, createdAt, updatedAt string

		err := rows.Scan(
			&appointment.ID,
			&appointment.PatientID,
			&appointment.DoctorID,
			&scheduledAt,
			&appointment.Duration,
			&appointment.Status,
			&appointment.Reason,
			&appointment.Notes,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Parse times
		appointment.ScheduledAt, _ = time.Parse(time.RFC3339, scheduledAt)
		appointment.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		appointment.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		appointments = append(appointments, &appointment)
	}

	return appointments, rows.Err()
}

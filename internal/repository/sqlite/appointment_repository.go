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
		INSERT INTO appointments (
			id, patient_id, doctor_id, scheduled_at, duration,
			reason, notes, status, created_at, updated_at,
			reminder_24h_sent, reminder_1h_sent, service_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		appointment.ID,
		appointment.PatientID,
		appointment.DoctorID,
		appointment.ScheduledAt,
		appointment.Duration,
		appointment.Reason,
		appointment.Notes,
		appointment.Status,
		appointment.CreatedAt,
		appointment.UpdatedAt,
		appointment.Reminder24hSent,
		appointment.Reminder1hSent,
		appointment.ServiceID,
	)

	return err
}

// FindByID retrieves an appointment by its unique identifier
func (r *SqliteAppointmentRepository) FindByID(ctx context.Context, id string) (*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, duration, status, reason, notes, created_at, updated_at, reminder_24h_sent, reminder_1h_sent
		FROM appointments
		WHERE id = $1
	`

	var appointment domain.Appointment
	var scheduledAt, createdAt, updatedAt time.Time

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
		&appointment.Reminder24hSent,
		&appointment.Reminder1hSent,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Assign scanned values
	appointment.ScheduledAt = scheduledAt
	appointment.CreatedAt = createdAt
	appointment.UpdatedAt = updatedAt

	return &appointment, nil
}

// FindByPatientID retrieves all appointments for a specific patient
func (r *SqliteAppointmentRepository) FindByPatientID(ctx context.Context, patientID string) ([]*domain.Appointment, error) {
	query := `
		SELECT
			a.id,
			a.patient_id,
			a.doctor_id,
			a.scheduled_at,
			a.duration,
			a.status,
			a.reason,
			a.notes,
			a.created_at,
			a.updated_at,
			a.reminder_24h_sent,
			a.reminder_1h_sent,
			(p.first_name || ' ' || p.last_name) as patient_name,
			(d.first_name || ' ' || d.last_name) as doctor_name
		FROM appointments a
		JOIN users p ON p.id = a.patient_id
		JOIN users d ON d.id = a.doctor_id
		WHERE a.patient_id = $1
		ORDER BY a.scheduled_at DESC
	`

	return r.queryAppointmentsWithNames(ctx, query, patientID)
}

// FindByDoctorID retrieves all appointments for a specific doctor
func (r *SqliteAppointmentRepository) FindByDoctorID(ctx context.Context, doctorID string) ([]*domain.Appointment, error) {
	query := `
		SELECT
			a.id,
			a.patient_id,
			a.doctor_id,
			a.scheduled_at,
			a.duration,
			a.status,
			a.reason,
			a.notes,
			a.created_at,
			a.updated_at,
			a.reminder_24h_sent,
			a.reminder_1h_sent,
			(p.first_name || ' ' || p.last_name) as patient_name,
			(d.first_name || ' ' || d.last_name) as doctor_name
		FROM appointments a
		JOIN users p ON p.id = a.patient_id
		JOIN users d ON d.id = a.doctor_id
		WHERE a.doctor_id = $1
		ORDER BY a.scheduled_at DESC
	`

	return r.queryAppointmentsWithNames(ctx, query, doctorID)
}

// FindByDoctorAndDate retrieves all appointments for a doctor on a specific date
func (r *SqliteAppointmentRepository) FindByDoctorAndDate(ctx context.Context, doctorID string, date time.Time) ([]*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, duration, status, reason, notes, created_at, updated_at, reminder_24h_sent, reminder_1h_sent
		FROM appointments
		WHERE doctor_id = $1 AND DATE(scheduled_at) = DATE($2)
		ORDER BY scheduled_at ASC
	`

	return r.queryAppointments(ctx, query, doctorID, date)
}

// FindByDoctorAndDateRange retrieves appointments for a doctor within a date range
func (r *SqliteAppointmentRepository) FindByDoctorAndDateRange(ctx context.Context, doctorID string, start, end time.Time) ([]*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, duration, status, reason, notes, created_at, updated_at, reminder_24h_sent, reminder_1h_sent
		FROM appointments
		WHERE doctor_id = $1 AND scheduled_at >= $2 AND scheduled_at < $3
		ORDER BY scheduled_at ASC
	`

	return r.queryAppointments(ctx, query, doctorID, start, end)
}

// Update modifies an existing appointment in the database
func (r *SqliteAppointmentRepository) Update(ctx context.Context, appointment *domain.Appointment) error {
	query := `
		UPDATE appointments
		SET status = $1, notes = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		appointment.Status,
		appointment.Notes,
		appointment.UpdatedAt,
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
	query := `DELETE FROM appointments WHERE id = $1`

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
		var scheduledAt, createdAt, updatedAt time.Time

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
			&appointment.Reminder24hSent,
			&appointment.Reminder1hSent,
		)

		if err != nil {
			return nil, err
		}

		// Assign scanned values
		appointment.ScheduledAt = scheduledAt
		appointment.CreatedAt = createdAt
		appointment.UpdatedAt = updatedAt

		appointments = append(appointments, &appointment)
	}

	return appointments, rows.Err()
}

// FindByScheduledAtRange finds appointments within a time range with specific status
func (r *SqliteAppointmentRepository) FindByScheduledAtRange(ctx context.Context, start, end time.Time, status string) ([]*domain.Appointment, error) {
	query := `
		SELECT id, patient_id, doctor_id, scheduled_at, duration, status, reason, notes, created_at, updated_at, reminder_24h_sent, reminder_1h_sent
		FROM appointments
		WHERE scheduled_at >= $1 AND scheduled_at <= $2 AND status = $3
		ORDER BY scheduled_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, start, end, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment

	for rows.Next() {
		var appointment domain.Appointment
		var scheduledAt, createdAt, updatedAt time.Time

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
			&appointment.Reminder24hSent,
			&appointment.Reminder1hSent,
		)

		if err != nil {
			return nil, err
		}

		appointment.ScheduledAt = scheduledAt
		appointment.CreatedAt = createdAt
		appointment.UpdatedAt = updatedAt

		appointments = append(appointments, &appointment)
	}

	return appointments, rows.Err()
}

// MarkReminder24hSent marks the 24-hour reminder as sent
func (r *SqliteAppointmentRepository) MarkReminder24hSent(ctx context.Context, id string) error {
	query := `UPDATE appointments SET reminder_24h_sent = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// MarkReminder1hSent marks the 1-hour reminder as sent
func (r *SqliteAppointmentRepository) MarkReminder1hSent(ctx context.Context, id string) error {
	query := `UPDATE appointments SET reminder_1h_sent = TRUE WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// queryAppointmentsWithNames is a helper method to query appointments with patient and doctor names
func (r *SqliteAppointmentRepository) queryAppointmentsWithNames(ctx context.Context, query string, args ...interface{}) ([]*domain.Appointment, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var appointments []*domain.Appointment

	for rows.Next() {
		var appointment domain.Appointment
		var scheduledAt, createdAt, updatedAt time.Time
		var patientName, doctorName string

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
			&appointment.Reminder24hSent,
			&appointment.Reminder1hSent,
			&patientName,
			&doctorName,
		)

		if err != nil {
			return nil, err
		}

		// Assign scanned values
		appointment.ScheduledAt = scheduledAt
		appointment.CreatedAt = createdAt
		appointment.UpdatedAt = updatedAt

		// Set names
		appointment.PatientName = patientName
		appointment.DoctorName = doctorName

		appointments = append(appointments, &appointment)
	}

	return appointments, rows.Err()
}

// CountByStatus counts appointments by status
func (r *SqliteAppointmentRepository) CountByStatus(ctx context.Context, status string) (int, error) {
	query := `SELECT COUNT(*) FROM appointments WHERE status = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, status).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountAll counts all appointments
func (r *SqliteAppointmentRepository) CountAll(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM appointments`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetTotalRevenue calculates total revenue from completed appointments
func (r *SqliteAppointmentRepository) GetTotalRevenue(ctx context.Context) (float64, error) {
	query := `
		SELECT COALESCE(SUM(s.price), 0)
		FROM appointments a
		JOIN services s ON a.service_id = s.id
		WHERE a.status = $1
	`

	var revenue float64
	err := r.db.QueryRowContext(ctx, query, "completed").Scan(&revenue)
	if err != nil {
		return 0, err
	}

	return revenue, nil
}

// GetRevenueByService gets revenue grouped by service
func (r *SqliteAppointmentRepository) GetRevenueByService(ctx context.Context) (map[string]struct {
	ServiceName string
	Count       int
	Revenue     float64
}, error) {
	query := `
		SELECT
			a.service_id,
			s.name,
			COUNT(*) as count,
			SUM(s.price) as revenue
		FROM appointments a
		JOIN services s ON a.service_id = s.id
		WHERE a.status = $1
		GROUP BY a.service_id, s.name
		ORDER BY revenue DESC
	`

	rows, err := r.db.QueryContext(ctx, query, "completed")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]struct {
		ServiceName string
		Count       int
		Revenue     float64
	})

	for rows.Next() {
		var serviceID, serviceName string
		var count int
		var revenue float64

		if err := rows.Scan(&serviceID, &serviceName, &count, &revenue); err != nil {
			return nil, err
		}

		result[serviceID] = struct {
			ServiceName string
			Count       int
			Revenue     float64
		}{
			ServiceName: serviceName,
			Count:       count,
			Revenue:     revenue,
		}
	}

	return result, rows.Err()
}

// GetTopDoctors gets doctors with most appointments
func (r *SqliteAppointmentRepository) GetTopDoctors(ctx context.Context, limit int) ([]struct {
	DoctorID              string
	TotalAppointments     int
	CompletedAppointments int
}, error) {
	query := `
		SELECT
			doctor_id,
			COUNT(*) as total,
			SUM(CASE WHEN status = 'completed' THEN 1 ELSE 0 END) as completed
		FROM appointments
		GROUP BY doctor_id
		ORDER BY total DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		DoctorID              string
		TotalAppointments     int
		CompletedAppointments int
	}

	for rows.Next() {
		var item struct {
			DoctorID              string
			TotalAppointments     int
			CompletedAppointments int
		}

		if err := rows.Scan(&item.DoctorID, &item.TotalAppointments, &item.CompletedAppointments); err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return results, rows.Err()
}

// GetTopServices gets most used services
func (r *SqliteAppointmentRepository) GetTopServices(ctx context.Context, limit int) ([]struct {
	ServiceID   string
	ServiceName string
	Count       int
}, error) {
	query := `
		SELECT
			a.service_id,
			s.name,
			COUNT(*) as count
		FROM appointments a
		JOIN services s ON a.service_id = s.id
		GROUP BY a.service_id, s.name
		ORDER BY count DESC
		LIMIT $1
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []struct {
		ServiceID   string
		ServiceName string
		Count       int
	}

	for rows.Next() {
		var item struct {
			ServiceID   string
			ServiceName string
			Count       int
		}

		if err := rows.Scan(&item.ServiceID, &item.ServiceName, &item.Count); err != nil {
			return nil, err
		}

		results = append(results, item)
	}

	return results, rows.Err()
}

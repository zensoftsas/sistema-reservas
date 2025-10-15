package sqlite

import (
	"context"
	"database/sql"
	"time"

	"version-1-0/internal/domain"
)

// SqliteScheduleRepository implements ScheduleRepository for SQLite
type SqliteScheduleRepository struct {
	db *sql.DB
}

// NewSqliteScheduleRepository creates a new SQLite schedule repository
func NewSqliteScheduleRepository(db *sql.DB) *SqliteScheduleRepository {
	return &SqliteScheduleRepository{db: db}
}

// Create creates a new schedule
func (r *SqliteScheduleRepository) Create(ctx context.Context, schedule *domain.Schedule) error {
	query := `
		INSERT INTO schedules (id, doctor_id, day_of_week, start_time, end_time, slot_duration, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		schedule.ID,
		schedule.DoctorID,
		schedule.DayOfWeek,
		schedule.StartTime,
		schedule.EndTime,
		schedule.SlotDuration,
		schedule.IsActive,
		schedule.CreatedAt,
		schedule.UpdatedAt,
	)

	return err
}

// FindByID finds a schedule by ID
func (r *SqliteScheduleRepository) FindByID(ctx context.Context, id string) (*domain.Schedule, error) {
	query := `
		SELECT id, doctor_id, day_of_week, start_time, end_time, slot_duration, is_active, created_at, updated_at
		FROM schedules
		WHERE id = $1
	`

	schedules, err := r.querySchedules(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(schedules) == 0 {
		return nil, nil
	}

	return schedules[0], nil
}

// FindByDoctorAndDay finds schedules for a doctor on a specific day
func (r *SqliteScheduleRepository) FindByDoctorAndDay(ctx context.Context, doctorID, dayOfWeek string) ([]*domain.Schedule, error) {
	query := `
		SELECT id, doctor_id, day_of_week, start_time, end_time, slot_duration, is_active, created_at, updated_at
		FROM schedules
		WHERE doctor_id = $1 AND day_of_week = $2 AND is_active = TRUE
		ORDER BY start_time ASC
	`

	return r.querySchedules(ctx, query, doctorID, dayOfWeek)
}

// FindByDoctor finds all schedules for a doctor
func (r *SqliteScheduleRepository) FindByDoctor(ctx context.Context, doctorID string) ([]*domain.Schedule, error) {
	query := `
		SELECT id, doctor_id, day_of_week, start_time, end_time, slot_duration, is_active, created_at, updated_at
		FROM schedules
		WHERE doctor_id = $1 AND is_active = TRUE
		ORDER BY
			CASE day_of_week
				WHEN 'monday' THEN 1
				WHEN 'tuesday' THEN 2
				WHEN 'wednesday' THEN 3
				WHEN 'thursday' THEN 4
				WHEN 'friday' THEN 5
				WHEN 'saturday' THEN 6
				WHEN 'sunday' THEN 7
			END,
			start_time ASC
	`

	return r.querySchedules(ctx, query, doctorID)
}

// Update updates a schedule
func (r *SqliteScheduleRepository) Update(ctx context.Context, schedule *domain.Schedule) error {
	query := `
		UPDATE schedules
		SET doctor_id = $1, day_of_week = $2, start_time = $3, end_time = $4,
		    slot_duration = $5, is_active = $6, updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		schedule.DoctorID,
		schedule.DayOfWeek,
		schedule.StartTime,
		schedule.EndTime,
		schedule.SlotDuration,
		schedule.IsActive,
		schedule.UpdatedAt,
		schedule.ID,
	)

	return err
}

// Delete deletes a schedule
func (r *SqliteScheduleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM schedules WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// DeleteByDoctorAndDay deletes all schedules for a doctor on a specific day
func (r *SqliteScheduleRepository) DeleteByDoctorAndDay(ctx context.Context, doctorID, dayOfWeek string) error {
	query := `DELETE FROM schedules WHERE doctor_id = $1 AND day_of_week = $2`
	_, err := r.db.ExecContext(ctx, query, doctorID, dayOfWeek)
	return err
}

// querySchedules is a helper method to query schedules
func (r *SqliteScheduleRepository) querySchedules(ctx context.Context, query string, args ...interface{}) ([]*domain.Schedule, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*domain.Schedule

	for rows.Next() {
		var schedule domain.Schedule
		var isActive bool
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&schedule.ID,
			&schedule.DoctorID,
			&schedule.DayOfWeek,
			&schedule.StartTime,
			&schedule.EndTime,
			&schedule.SlotDuration,
			&isActive,
			&createdAt,
			&updatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Assign scanned values
		schedule.IsActive = isActive
		schedule.CreatedAt = createdAt
		schedule.UpdatedAt = updatedAt

		schedules = append(schedules, &schedule)
	}

	return schedules, rows.Err()
}

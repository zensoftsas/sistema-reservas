package sqlite

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a database migration with version control
type Migration struct {
	Version     int
	Description string
	Up          func(*sql.DB) error
}

// migrations contains all database migrations in order
var migrations = []Migration{
	{
		Version:     1,
		Description: "Create initial schema",
		Up:          migrateV1_InitialSchema,
	},
	{
		Version:     2,
		Description: "Add reminder fields to appointments",
		Up:          migrateV2_AddReminderFields,
	},
	{
		Version:     3,
		Description: "Add service_id to appointments",
		Up:          migrateV3_AddServiceID,
	},
}

// runMigrations executes all pending migrations
func runMigrations(db *sql.DB) error {
	// Create schema_migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get current version
	currentVersion, err := getCurrentVersion(db)
	if err != nil {
		return fmt.Errorf("failed to get current version: %v", err)
	}

	log.Printf("Current database version: %d", currentVersion)

	// Run pending migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			log.Printf("Running migration v%d: %s", migration.Version, migration.Description)

			// Execute migration
			if err := migration.Up(db); err != nil {
				return fmt.Errorf("migration v%d failed: %v", migration.Version, err)
			}

			// Update version
			if err := setVersion(db, migration.Version); err != nil {
				return fmt.Errorf("failed to update version to %d: %v", migration.Version, err)
			}

			log.Printf("✓ Migration v%d completed successfully", migration.Version)
		}
	}

	return nil
}

// createMigrationsTable creates the schema_migrations table
func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

// getCurrentVersion returns the current database version
func getCurrentVersion(db *sql.DB) (int, error) {
	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

// setVersion records a migration version as applied
func setVersion(db *sql.DB, version int) error {
	_, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
	return err
}

// migrateV1_InitialSchema creates all initial tables
func migrateV1_InitialSchema(db *sql.DB) error {
	// Create users table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			first_name TEXT NOT NULL,
			last_name TEXT NOT NULL,
			phone TEXT NOT NULL,
			role TEXT NOT NULL CHECK(role IN ('admin', 'doctor', 'patient')),
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`); err != nil {
		return err
	}

	// Create indexes on users
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`); err != nil {
		return err
	}

	// Create patients table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS patients (
			id TEXT PRIMARY KEY,
			user_id TEXT UNIQUE NOT NULL,
			birthdate TIMESTAMP NOT NULL,
			document_type TEXT NOT NULL,
			document_number TEXT NOT NULL,
			address TEXT NOT NULL,
			emergency_contact_name TEXT NOT NULL,
			emergency_contact_phone TEXT NOT NULL,
			blood_type TEXT,
			allergies TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_patients_user_id ON patients(user_id)`); err != nil {
		return err
	}

	// Create doctors table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS doctors (
			id TEXT PRIMARY KEY,
			user_id TEXT UNIQUE NOT NULL,
			specialty TEXT NOT NULL,
			license_number TEXT UNIQUE NOT NULL,
			years_of_experience INTEGER NOT NULL,
			education TEXT,
			bio TEXT,
			consultation_fee REAL NOT NULL,
			is_available BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_doctors_user_id ON doctors(user_id)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_doctors_specialty ON doctors(specialty)`); err != nil {
		return err
	}

	// Create schedules table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schedules (
			id TEXT PRIMARY KEY,
			doctor_id TEXT NOT NULL,
			day_of_week TEXT NOT NULL CHECK(day_of_week IN ('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday')),
			start_time TEXT NOT NULL,
			end_time TEXT NOT NULL,
			slot_duration INTEGER NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE
		)
	`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_schedules_doctor_id ON schedules(doctor_id)`); err != nil {
		return err
	}

	// Create services table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS services (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			duration_minutes INTEGER NOT NULL CHECK(duration_minutes > 0),
			price REAL NOT NULL CHECK(price >= 0),
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_services_name ON services(name)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_services_is_active ON services(is_active)`); err != nil {
		return err
	}

	// Create doctor_services table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS doctor_services (
			id TEXT PRIMARY KEY,
			doctor_id TEXT NOT NULL,
			service_id TEXT NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE,
			FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE,
			UNIQUE(doctor_id, service_id)
		)
	`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_doctor_services_doctor_id ON doctor_services(doctor_id)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_doctor_services_service_id ON doctor_services(service_id)`); err != nil {
		return err
	}

	// Create appointments table
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS appointments (
			id TEXT PRIMARY KEY,
			patient_id TEXT NOT NULL,
			doctor_id TEXT NOT NULL,
			scheduled_at TIMESTAMP NOT NULL,
			duration INTEGER NOT NULL,
			reason TEXT NOT NULL,
			notes TEXT,
			status TEXT NOT NULL CHECK(status IN ('pending', 'confirmed', 'cancelled', 'completed')),
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			cancelled_at TIMESTAMP,
			cancellation_reason TEXT,
			FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE,
			FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE
		)
	`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_appointments_patient_id ON appointments(patient_id)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id ON appointments(doctor_id)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_appointments_scheduled_at ON appointments(scheduled_at)`); err != nil {
		return err
	}
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status)`); err != nil {
		return err
	}

	return nil
}

// migrateV2_AddReminderFields adds reminder tracking fields to appointments
func migrateV2_AddReminderFields(db *sql.DB) error {
	// Check if columns exist before adding
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_name='appointments' AND column_name='reminder_24h_sent'
	`).Scan(&count)

	if err != nil || count == 0 {
		if _, err := db.Exec(`ALTER TABLE appointments ADD COLUMN reminder_24h_sent BOOLEAN DEFAULT FALSE`); err != nil {
			return err
		}
	}

	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_name='appointments' AND column_name='reminder_1h_sent'
	`).Scan(&count)

	if err != nil || count == 0 {
		if _, err := db.Exec(`ALTER TABLE appointments ADD COLUMN reminder_1h_sent BOOLEAN DEFAULT FALSE`); err != nil {
			return err
		}
	}

	return nil
}

// migrateV3_AddServiceID adds service_id foreign key to appointments
func migrateV3_AddServiceID(db *sql.DB) error {
	// Check if column exists before adding
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_name='appointments' AND column_name='service_id'
	`).Scan(&count)

	if err != nil || count == 0 {
		if _, err := db.Exec(`ALTER TABLE appointments ADD COLUMN service_id TEXT REFERENCES services(id)`); err != nil {
			return err
		}
		if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_appointments_service_id ON appointments(service_id)`); err != nil {
			return err
		}
	}

	return nil
}

package sqlite

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

// InitDB initializes and returns a SQLite database connection
// Returns an error if the connection cannot be established
func InitDB(filepath string) (*sql.DB, error) {
	// Open database connection
	db, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	// Verify connection is working
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	// Configure SQLite connection pool
	// SQLite works best with a single connection to avoid locking issues
	db.SetMaxOpenConns(1)

	log.Printf("Database connection established: %s", filepath)

	return db, nil
}

// MigrateDB creates all necessary database tables and indexes
// Returns an error if migration fails
func MigrateDB(db *sql.DB) error {
	// Create users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		phone TEXT NOT NULL,
		role TEXT NOT NULL CHECK(role IN ('admin', 'doctor', 'patient')),
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	if _, err := db.Exec(createUsersTable); err != nil {
		return err
	}

	// Create index on email for faster lookups
	createEmailIndex := `CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);`
	if _, err := db.Exec(createEmailIndex); err != nil {
		return err
	}

	// Create index on role for filtering
	createRoleIndex := `CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);`
	if _, err := db.Exec(createRoleIndex); err != nil {
		return err
	}

	// Create patients table
	createPatientsTable := `
	CREATE TABLE IF NOT EXISTS patients (
		id TEXT PRIMARY KEY,
		user_id TEXT UNIQUE NOT NULL,
		birthdate DATETIME NOT NULL,
		document_type TEXT NOT NULL,
		document_number TEXT NOT NULL,
		address TEXT NOT NULL,
		emergency_contact_name TEXT NOT NULL,
		emergency_contact_phone TEXT NOT NULL,
		blood_type TEXT,
		allergies TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(createPatientsTable); err != nil {
		return err
	}

	// Create index on user_id for patients
	createPatientUserIndex := `CREATE INDEX IF NOT EXISTS idx_patients_user_id ON patients(user_id);`
	if _, err := db.Exec(createPatientUserIndex); err != nil {
		return err
	}

	// Create doctors table
	createDoctorsTable := `
	CREATE TABLE IF NOT EXISTS doctors (
		id TEXT PRIMARY KEY,
		user_id TEXT UNIQUE NOT NULL,
		specialty TEXT NOT NULL,
		license_number TEXT UNIQUE NOT NULL,
		years_of_experience INTEGER NOT NULL,
		education TEXT,
		bio TEXT,
		consultation_fee REAL NOT NULL,
		is_available BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(createDoctorsTable); err != nil {
		return err
	}

	// Create indexes for doctors
	createDoctorUserIndex := `CREATE INDEX IF NOT EXISTS idx_doctors_user_id ON doctors(user_id);`
	if _, err := db.Exec(createDoctorUserIndex); err != nil {
		return err
	}

	createDoctorSpecialtyIndex := `CREATE INDEX IF NOT EXISTS idx_doctors_specialty ON doctors(specialty);`
	if _, err := db.Exec(createDoctorSpecialtyIndex); err != nil {
		return err
	}

	// Create schedules table
	createSchedulesTable := `
	CREATE TABLE IF NOT EXISTS schedules (
		id TEXT PRIMARY KEY,
		doctor_id TEXT NOT NULL,
		day_of_week TEXT NOT NULL CHECK(day_of_week IN ('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday')),
		start_time TEXT NOT NULL,
		end_time TEXT NOT NULL,
		slot_duration INTEGER NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(createSchedulesTable); err != nil {
		return err
	}

	// Create index on doctor_id for schedules
	createScheduleDoctorIndex := `CREATE INDEX IF NOT EXISTS idx_schedules_doctor_id ON schedules(doctor_id);`
	if _, err := db.Exec(createScheduleDoctorIndex); err != nil {
		return err
	}

	// Create services table
	createServicesTable := `
	CREATE TABLE IF NOT EXISTS services (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		duration_minutes INTEGER NOT NULL CHECK(duration_minutes > 0),
		price REAL NOT NULL CHECK(price >= 0),
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	if _, err := db.Exec(createServicesTable); err != nil {
		return err
	}

	// Create index on service name for faster lookups
	createServiceNameIndex := `CREATE INDEX IF NOT EXISTS idx_services_name ON services(name);`
	if _, err := db.Exec(createServiceNameIndex); err != nil {
		return err
	}

	// Create index on is_active for filtering active services
	createServiceActiveIndex := `CREATE INDEX IF NOT EXISTS idx_services_is_active ON services(is_active);`
	if _, err := db.Exec(createServiceActiveIndex); err != nil {
		return err
	}

	// Create doctor_services table (many-to-many relationship)
	createDoctorServicesTable := `
	CREATE TABLE IF NOT EXISTS doctor_services (
		id TEXT PRIMARY KEY,
		doctor_id TEXT NOT NULL,
		service_id TEXT NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE,
		FOREIGN KEY (service_id) REFERENCES services(id) ON DELETE CASCADE,
		UNIQUE(doctor_id, service_id)
	);`

	if _, err := db.Exec(createDoctorServicesTable); err != nil {
		return err
	}

	// Create indexes for doctor_services
	createDoctorServicesDocIndex := `CREATE INDEX IF NOT EXISTS idx_doctor_services_doctor_id ON doctor_services(doctor_id);`
	if _, err := db.Exec(createDoctorServicesDocIndex); err != nil {
		return err
	}

	createDoctorServicesServiceIndex := `CREATE INDEX IF NOT EXISTS idx_doctor_services_service_id ON doctor_services(service_id);`
	if _, err := db.Exec(createDoctorServicesServiceIndex); err != nil {
		return err
	}

	// Create appointments table
	createAppointmentsTable := `
	CREATE TABLE IF NOT EXISTS appointments (
		id TEXT PRIMARY KEY,
		patient_id TEXT NOT NULL,
		doctor_id TEXT NOT NULL,
		scheduled_at DATETIME NOT NULL,
		duration INTEGER NOT NULL,
		reason TEXT NOT NULL,
		notes TEXT,
		status TEXT NOT NULL CHECK(status IN ('pending', 'confirmed', 'cancelled', 'completed')),
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		cancelled_at DATETIME,
		cancellation_reason TEXT,
		FOREIGN KEY (patient_id) REFERENCES patients(id) ON DELETE CASCADE,
		FOREIGN KEY (doctor_id) REFERENCES doctors(id) ON DELETE CASCADE
	);`

	if _, err := db.Exec(createAppointmentsTable); err != nil {
		return err
	}

	// Add reminder fields to appointments table if they don't exist
	alterAppointmentsReminder24h := `
	ALTER TABLE appointments ADD COLUMN reminder_24h_sent BOOLEAN DEFAULT 0;
`
	db.Exec(alterAppointmentsReminder24h) // Ignore error if column exists

	alterAppointmentsReminder1h := `
	ALTER TABLE appointments ADD COLUMN reminder_1h_sent BOOLEAN DEFAULT 0;
`
	db.Exec(alterAppointmentsReminder1h) // Ignore error if column exists

	// Add service_id to appointments table if it doesn't exist
	alterAppointmentsServiceID := `
	ALTER TABLE appointments ADD COLUMN service_id TEXT REFERENCES services(id);
`
	db.Exec(alterAppointmentsServiceID) // Ignore error if column exists

	// Create indexes for appointments
	createAppointmentPatientIndex := `CREATE INDEX IF NOT EXISTS idx_appointments_patient_id ON appointments(patient_id);`
	if _, err := db.Exec(createAppointmentPatientIndex); err != nil {
		return err
	}

	createAppointmentDoctorIndex := `CREATE INDEX IF NOT EXISTS idx_appointments_doctor_id ON appointments(doctor_id);`
	if _, err := db.Exec(createAppointmentDoctorIndex); err != nil {
		return err
	}

	createAppointmentScheduledIndex := `CREATE INDEX IF NOT EXISTS idx_appointments_scheduled_at ON appointments(scheduled_at);`
	if _, err := db.Exec(createAppointmentScheduledIndex); err != nil {
		return err
	}

	createAppointmentStatusIndex := `CREATE INDEX IF NOT EXISTS idx_appointments_status ON appointments(status);`
	if _, err := db.Exec(createAppointmentStatusIndex); err != nil {
		return err
	}

	createAppointmentServiceIndex := `CREATE INDEX IF NOT EXISTS idx_appointments_service_id ON appointments(service_id);`
	if _, err := db.Exec(createAppointmentServiceIndex); err != nil {
		return err
	}

	log.Println("Database migration completed successfully")

	return nil
}

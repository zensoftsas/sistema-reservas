package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

func main() {
	// Open database
	db, err := sql.Open("sqlite", "clinica.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check if doctor profile exists for dr.garcia
	var count int
	query := `
		SELECT COUNT(*)
		FROM doctors d
		INNER JOIN users u ON d.user_id = u.id
		WHERE u.email = 'dr.garcia@clinica.com'
	`
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count > 0 {
		fmt.Println("✅ El doctor ya tiene perfil en la tabla doctors")

		// Show existing doctor info
		var specialty string
		query = `
			SELECT d.specialty
			FROM doctors d
			INNER JOIN users u ON d.user_id = u.id
			WHERE u.email = 'dr.garcia@clinica.com'
		`
		db.QueryRow(query).Scan(&specialty)
		fmt.Printf("   Especialidad: %s\n", specialty)
		return
	}

	fmt.Println("⚠️  El doctor NO tiene perfil en tabla doctors. Creando...")

	// Get doctor user_id
	var doctorUserID string
	err = db.QueryRow("SELECT id FROM users WHERE email = 'dr.garcia@clinica.com'").Scan(&doctorUserID)
	if err != nil {
		log.Fatal("Error: Doctor user not found. Create user first.")
	}

	// Insert doctor profile
	insertQuery := `
		INSERT INTO doctors (
			id, user_id, specialty, license_number, years_of_experience,
			education, bio, consultation_fee, is_available, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	doctorID := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	_, err = db.Exec(
		insertQuery,
		doctorID,
		doctorUserID,
		"Cardiología", // specialty
		"MED-2024-001", // license_number
		10, // years_of_experience
		"Universidad Nacional Mayor de San Marcos - Especialización en Cardiología", // education
		"Especialista en cardiología con 10 años de experiencia en tratamiento de enfermedades cardiovasculares", // bio
		150.00, // consultation_fee
		true, // is_available
		now,
		now,
	)

	if err != nil {
		log.Fatal("Error inserting doctor profile:", err)
	}

	fmt.Println("✅ Perfil de doctor creado exitosamente!")
	fmt.Println("   ID:", doctorID)
	fmt.Println("   Especialidad: Cardiología")
	fmt.Println("   Años de experiencia: 10")
	fmt.Println("   Tarifa: S/. 150.00")
}
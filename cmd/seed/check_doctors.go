package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "clinica.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check users with role doctor
	var userDoctorCount int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE role='doctor'").Scan(&userDoctorCount)
	fmt.Printf("Users with role doctor: %d\n", userDoctorCount)

	// Check doctor profiles
	var doctorProfileCount int
	db.QueryRow("SELECT COUNT(*) FROM doctors").Scan(&doctorProfileCount)
	fmt.Printf("Profiles in doctors table: %d\n\n", doctorProfileCount)

	// Show doctors with profile
	fmt.Println("Doctors WITH complete profile:")
	rows, _ := db.Query(`
		SELECT u.email, u.first_name, u.last_name, d.specialty
		FROM users u
		INNER JOIN doctors d ON u.id = d.user_id
		WHERE u.role='doctor'
	`)
	defer rows.Close()

	count := 0
	for rows.Next() {
		var email, firstName, lastName, specialty string
		rows.Scan(&email, &firstName, &lastName, &specialty)
		fmt.Printf("   - %s %s (%s) - %s\n", firstName, lastName, email, specialty)
		count++
	}
	if count == 0 {
		fmt.Println("   (none)")
	}

	// Show doctors WITHOUT profile
	fmt.Println("\nDoctors WITHOUT profile in doctors table:")
	rows2, _ := db.Query(`
		SELECT u.email, u.first_name, u.last_name
		FROM users u
		LEFT JOIN doctors d ON u.id = d.user_id
		WHERE u.role='doctor' AND d.id IS NULL
	`)
	defer rows2.Close()

	count2 := 0
	for rows2.Next() {
		var email, firstName, lastName string
		rows2.Scan(&email, &firstName, &lastName)
		fmt.Printf("   - %s %s (%s)\n", firstName, lastName, email)
		count2++
	}
	if count2 == 0 {
		fmt.Println("   (none)")
	}
}

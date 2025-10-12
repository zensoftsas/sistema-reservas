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

	// Update specialty to remove accents
	_, err = db.Exec("UPDATE doctors SET specialty = 'Cardiologia' WHERE specialty = 'Cardiología'")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Specialty normalized: Cardiología -> Cardiologia")

	// Verify
	var specialty string
	db.QueryRow("SELECT specialty FROM doctors WHERE user_id IN (SELECT id FROM users WHERE email = 'dr.garcia@clinica.com')").Scan(&specialty)
	fmt.Printf("New specialty: %s\n", specialty)
}

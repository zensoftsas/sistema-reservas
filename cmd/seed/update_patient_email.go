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

	// Update patient email to real email
	newEmail := "zensoftsas@gmail.com"

	result, err := db.Exec("UPDATE users SET email = ? WHERE email = 'paciente@clinica.com'", newEmail)
	if err != nil {
		log.Fatal("Error updating email:", err)
	}

	rows, _ := result.RowsAffected()
	if rows > 0 {
		fmt.Printf("Email actualizado: paciente@clinica.com → %s\n", newEmail)
	} else {
		fmt.Println("No se encontró el paciente")
	}
}

package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// InitDB initializes and returns a PostgreSQL database connection
// Returns an error if the connection cannot be established
func InitDB(databaseURL string) (*sql.DB, error) {
	// Open PostgreSQL connection
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Printf("Database connection established: PostgreSQL (Neon)")

	// Run migrations with version control
	if err := runMigrations(db); err != nil {
		return nil, err
	}

	log.Println("Database migrations completed successfully")

	return db, nil
}

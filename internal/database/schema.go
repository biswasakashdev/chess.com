package database

import (
	"database/sql"
	"log"
	"os"
)

func InitSchema(db *sql.DB) error {

	// Read SQL file
	script, err := os.ReadFile("scripts/sql/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	// Execute SQL
	_, err = db.Exec(string(script))
	if err != nil {
		log.Fatal(err)
	}

	return err
}

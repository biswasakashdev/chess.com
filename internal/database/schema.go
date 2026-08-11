package database

import "database/sql"

func InitSchema(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		username TEXT NOT NULL UNIQUE,
		hashed_password TEXT NOT NULL,
		created_at DATETIME NOT NULL
	);
	`

	_, err := db.Exec(query)

	return err
}

package tables

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(up, down)
}

func up(tx *sql.Tx) error {

	if _, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS profiles (
			id INT PRIMARY KEY,
			name VARCHAR(256),
			age INT,
			description VARCHAR(256),
			user_location VARCHAR(256),
			location TEXT ARRAY
		)
	`); err != nil {
		return err
	}

	return nil
}

func down(tx *sql.Tx) error {

	if _, err := tx.Exec(`DROP TABLE profiles`); err != nil {
		return err
	}

	return nil
}

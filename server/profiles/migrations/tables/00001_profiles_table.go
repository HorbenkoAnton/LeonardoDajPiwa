package tables

import (
	"database/sql"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(upProfilesTable, downProfilesTable)
}

func upProfilesTable(tx *sql.Tx) error {

	if _, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS profiles (
			id INT PRIMARY KEY,
			name VARCHAR(256),
			age INT,
			description VARCHAR(256),
		    pfp_id VARCHAR(256),
			user_location VARCHAR(256),
			location TEXT ARRAY
		)
	`); err != nil {
		return err
	}

	return nil
}

func downProfilesTable(tx *sql.Tx) error {

	if _, err := tx.Exec(`DROP TABLE profiles`); err != nil {
		return err
	}

	return nil
}

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
			id BIGINT PRIMARY KEY,
			name TEXT,
			age INT,
			description TEXT,
		    pfp_id TEXT,
			user_location TEXT,
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

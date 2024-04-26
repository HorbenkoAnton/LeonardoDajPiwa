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
		CREATE TABLE IF NOT EXISTS likes (
			liker INT NOT NULL,
			liked INT NOT NULL,
		    UNIQUE (liker, liked)
		)
	`); err != nil {
		return err
	}

	return nil
}

func down(tx *sql.Tx) error {

	if _, err := tx.Exec(`DROP TABLE likes`); err != nil {
		return err
	}

	return nil
}

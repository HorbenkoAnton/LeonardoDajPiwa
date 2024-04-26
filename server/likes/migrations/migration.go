package migrations

import (
	"database/sql"
	"github.com/pressly/goose"
	_ "likes/migrations/tables"
)

func Migrate(reload bool, db *sql.DB) {
	if reload {
		if err := goose.DownTo(db, ".", 0); err != nil {
			panic(err)
		}
	}
	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}
}

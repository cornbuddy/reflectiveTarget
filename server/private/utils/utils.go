package utils

import (
	"database/sql"
	_ "embed"
)

//go:embed tables.sql
var initQuery string

func InitDatabase(db *sql.DB) error {
	if _, err := db.Exec(initQuery); err != nil {
		return err
	}

	return nil
}

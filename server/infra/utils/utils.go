package utils

import (
	"database/sql"

	"github.com/cornbuddy/reflectiveTarget/server/infra/migrations"
)

func InitDatabase(db *sql.DB) error {
	if _, err := db.Exec(migrations.InitQuery); err != nil {
		return err
	}

	return nil
}

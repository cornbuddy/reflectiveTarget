package utils

import (
	"context"
	"database/sql"

	"github.com/cornbuddy/reflectiveTarget/server/infra/migrations"
)

func InitDatabase(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, migrations.InitQuery); err != nil {
		return err
	}

	return nil
}

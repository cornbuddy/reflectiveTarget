package utils

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/cornbuddy/reflectiveTarget/server/infra/migrations"
)

const DbImage = "postgres:18-alpine"

// starts and initializes database
func SetupDB(ctx context.Context) (Cleanup, *sql.DB, error) {
	cleanup, db, err := StartDB(ctx)
	if err != nil {
		return cleanup, nil, err
	}

	if _, err := db.ExecContext(ctx, migrations.InitQuery); err != nil {
		return cleanup, nil, err
	}

	return cleanup, db, nil
}

// starts testcontainer with database
func StartDB(ctx context.Context) (Cleanup, *sql.DB, error) {
	//nolint:mnd
	str := wait.ForLog("database system is ready to accept connections").
		WithOccurrence(2).
		WithStartupTimeout(5 * time.Second)
	opts := []tc.ContainerCustomizer{
		tc.WithWaitStrategy(str),
	}
	cont, err := postgres.Run(ctx, DbImage, opts...)

	cleanup := func() error {
		if err := tc.TerminateContainer(cont); err != nil {
			return err
		}

		return nil
	}

	if err != nil {
		return cleanup, nil, err
	}

	connStr, err := cont.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return cleanup, nil, err
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return cleanup, nil, err
	}

	cleanup = func() error {
		if err := db.Close(); err != nil {
			return err
		}

		if err := tc.TerminateContainer(cont); err != nil {
			return err
		}

		return nil
	}

	return cleanup, db, nil
}

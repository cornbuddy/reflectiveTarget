package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const Image = "postgres:18-alpine"

type Cleanup func() error

func SetupTestDb(ctx context.Context, t *testing.T) (Cleanup, *sql.DB, error) {
	t.Helper()

	emptyCleanup := func() error { return nil }

	str := wait.ForLog("database system is ready to accept connections").
		WithOccurrence(2).
		WithStartupTimeout(5 * time.Second)
	opts := []tc.ContainerCustomizer{
		tc.WithWaitStrategy(str),
	}
	dbContainer, err := postgres.Run(ctx, Image, opts...)
	if err != nil {
		return emptyCleanup, nil, err
	}

	connStr, err := dbContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return emptyCleanup, nil, err
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return emptyCleanup, nil, err
	}

	cleanup := func() error {
		if err = db.Close(); err != nil {
			return err
		}

		if err := tc.TerminateContainer(dbContainer); err != nil {
			return err
		}

		return nil
	}
	return cleanup, db, nil
}

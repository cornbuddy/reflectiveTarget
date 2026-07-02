package postgres

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
	"github.com/cornbuddy/reflectiveTarget/server/infra/migrations"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(ctx context.Context, db *sql.DB) (*PostgresStore, error) {
	log := log.Logger(ctx, zap.String("session-store", "postgres"))
	log.Debug("going to initialize store")
	_, err := db.ExecContext(ctx, migrations.InitSessionStoreQuery)
	if err != nil {
		return nil, err
	}

	log.Info("store is initialized")

	return &PostgresStore{db}, nil
}

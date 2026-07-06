package session

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"

	appsession "github.com/cornbuddy/reflectiveTarget/server/app/session"
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

func (s *PostgresStore) Update(
	ctx context.Context, id appsession.SessionID, data appsession.Data,
) error {
	return errors.New("not implemented")
}

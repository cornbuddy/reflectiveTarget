package db

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const Image = "postgres:18-alpine"

func Setup(ctx context.Context) (*postgres.PostgresContainer, error) {
	db, err := postgres.Run(ctx, Image)
	if err != nil {
		return nil, err
	}

	return db, nil
}

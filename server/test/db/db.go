package db

import (
	"context"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const Image = "postgres:18-alpine"

func Setup(ctx context.Context) (*postgres.PostgresContainer, error) {
	str := wait.ForLog("database system is ready to accept connections").
		WithOccurrence(2).
		WithStartupTimeout(5 * time.Second)
	opts := []testcontainers.ContainerCustomizer{
		testcontainers.WithWaitStrategy(str),
	}

	db, err := postgres.Run(ctx, Image, opts...)
	if err != nil {
		return nil, err
	}

	return db, nil
}

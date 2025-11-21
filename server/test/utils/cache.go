package utils

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	tc "github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

const CacheImage = "redis:8"

func SetupCache(ctx context.Context) (Cleanup, *redis.Client, error) {
	cont, err := tcredis.Run(ctx, CacheImage)

	cleanup := func() error {
		if err := tc.TerminateContainer(cont); err != nil {
			return err
		}

		return nil
	}

	if err != nil {
		return cleanup, nil, err
	}

	uri, err := cont.ConnectionString(ctx)
	if err != nil {
		return cleanup, nil, err
	}

	opts, err := redis.ParseURL(uri)
	if err != nil {
		return cleanup, nil, err
	}

	client := redis.NewClient(opts)
	if client == nil {
		return cleanup, nil, errors.New("failed to create redis client")
	}

	if err := client.Ping(ctx).Err(); err != nil {
		return cleanup, nil, err
	}

	cleanup = func() error {
		if err := client.Close(); err != nil {
			return err
		}

		if err := tc.TerminateContainer(cont); err != nil {
			return err
		}

		return nil
	}

	return cleanup, client, nil
}

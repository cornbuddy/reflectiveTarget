package daos

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/cornbuddy/reflectiveTarget/server/domain/constants"
)

type SessionStore struct {
	Ctx   context.Context
	Cache *redis.Client
}

const SessionPrefix = "session"

func (s SessionStore) GetUsernameFromSession(token string) (string, error) {
	ctx := s.Ctx
	cache := s.Cache
	key := fmt.Sprintf("%s:%s", SessionPrefix, token)

	username, err := cache.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return username, nil
}

func (s SessionStore) SaveSessionForUser(username, token string) error {
	ctx := s.Ctx
	cache := s.Cache
	expiration := constants.SessionDuration
	key := fmt.Sprintf("%s:%s", SessionPrefix, token)

	_, err := cache.Set(ctx, key, username, expiration).Result()
	if err != nil {
		return err
	}

	return nil
}

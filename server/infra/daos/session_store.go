package daos

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/cornbuddy/reflectiveTarget/server/domain/constants"
)

type SessionStore struct {
	Cache *redis.Client
}

func (s SessionStore) IsAuthenticated(
	ctx context.Context, token string,
) (*bool, error) {

	key := s.isAuthenticatedKey(token)
	value, err := s.Cache.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	isAuthenticated, err := strconv.ParseBool(value)
	if err != nil {
		return nil, err
	}

	return &isAuthenticated, nil
}

func (s SessionStore) SaveSession(
	ctx context.Context, token string, isAuthenticated bool,
) error {

	expiration := constants.SessionDuration
	key := s.isAuthenticatedKey(token)
	value := strconv.FormatBool(isAuthenticated)
	_, err := s.Cache.Set(ctx, key, value, expiration).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s SessionStore) isAuthenticatedKey(token string) string {
	const prefix = "session"
	const postfix = "isAuthenticated"

	return fmt.Sprintf("%s:%s:%s", prefix, token, postfix)
}

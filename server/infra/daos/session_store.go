package daos

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/domain/constants"
)

type SessionStore struct {
	Cache *redis.Client
}

func (s SessionStore) Get(
	ctx context.Context, sessionId string,
) (*sessiondata.SessionData, error) {

	key := s.isAuthenticatedKey(sessionId)
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

	data := sessiondata.SessionData{
		IsAuthenticated: isAuthenticated,
	}

	return &data, nil
}

func (s SessionStore) Update(
	ctx context.Context, sessionId string, data sessiondata.SessionData,
) error {

	expiration := constants.SessionDuration
	key := s.isAuthenticatedKey(sessionId)
	value := strconv.FormatBool(data.IsAuthenticated)
	_, err := s.Cache.Set(ctx, key, value, expiration).Result()
	if err != nil {
		return err
	}

	return nil
}

func (s SessionStore) isAuthenticatedKey(sessionId string) string {
	const prefix = "session"
	const postfix = "isAuthenticated"

	return fmt.Sprintf("%s:%s:%s", prefix, sessionId, postfix)
}

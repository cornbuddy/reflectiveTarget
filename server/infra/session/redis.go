package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/cornbuddy/reflectiveTarget/server/app/session"
)

const (
	prefix = "session"
)

type RedisStore struct {
	Cache *redis.Client
}

func (s RedisStore) Get(
	ctx context.Context, sessionId session.SessionID,
) (*session.Data, error) {
	key := s.Key(sessionId)
	jsonData, err := s.Cache.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var data session.Data
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (s RedisStore) Update(
	ctx context.Context, sessionId session.SessionID,
	data session.Data,
) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	key := s.Key(sessionId)
	err = s.Cache.Set(ctx, key, jsonData, session.Duration).Err()
	if err != nil {
		return err
	}

	return nil
}

// returns key to the session data
func (s RedisStore) Key(sessionId session.SessionID) string {
	return fmt.Sprintf("%s:%s", prefix, sessionId)
}

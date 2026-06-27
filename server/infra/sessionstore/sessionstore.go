package sessionstore

import (
	"context"

	redispkg "github.com/redis/go-redis/v9"

	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/infra/sessionstore/private/redis"
)

type SessionStore interface {
	Get(context.Context, sessiondata.SessionID) (*sessiondata.SessionData, error)
	Update(context.Context, sessiondata.SessionID, sessiondata.SessionData) error
}

func NewRedisStore(client *redispkg.Client) SessionStore {
	return &redis.RedisStore{Cache: client}
}

package session_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appsession "github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/infra/session"
)

func TestRedisStore(t *testing.T) {
	t.Parallel()

	store := session.RedisStore{cache}

	id := appsession.MakeID()
	want := appsession.Data{
		UserID:   69,
		Username: "kek",
	}
	require.NoError(t, store.Update(ctx, id, want))

	got, err := store.Get(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, want, *got, "operations should be deterministic")

	key := store.Key(id)
	ttl, err := store.Cache.TTL(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, appsession.Duration, ttl, "should have proper ttl")
}

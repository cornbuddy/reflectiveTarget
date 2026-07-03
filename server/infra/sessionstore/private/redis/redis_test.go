package redis_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
)

func TestShouldSetProperTTL(t *testing.T) {
	t.Parallel()

	id := sessiondata.SessionID("kek")
	data := sessiondata.SessionData{
		IsAuthenticated: true,
		UserID:          69,
		Username:        "kek",
	}
	require.NoError(t, store.Update(ctx, id, data))

	key := store.Key(id)
	ttl, err := store.Cache.TTL(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, sessiondata.SessionDuration, ttl)
}

package daos

import (
	"context"
	"fmt"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/domain/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionStoreShouldFetchUsername(t *testing.T) {
	t.Parallel()

	username := "username"
	token := "kek-token"
	require.NoError(t, store.SaveSessionForUser(username, token))

	type testCase struct {
		desc  string
		token string
		want  string
	}

	testCases := []testCase{{
		desc:  "should return nil when no session exists",
		token: "not exists",
		want:  "",
	}, {
		desc:  "should return username when session exists",
		token: token,
		want:  username,
	}}

	for _, tc := range testCases {
		got, err := store.GetUsernameFromSession(tc.token)
		require.NoError(t, err)
		assert.Equal(t, tc.want, got)
	}
}

func TestSessionStoreShouldSaveSessionForUser(t *testing.T) {
	t.Parallel()

	username := "kek"
	token := "kek"
	require.NoError(t, store.SaveSessionForUser(username, token))

	ctx := context.TODO()
	key := fmt.Sprintf("%s:%s", SessionPrefix, token)
	got, err := cache.Get(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, username, got)

	ttl, err := cache.TTL(ctx, key).Result()
	require.NoError(t, err)
	assert.Equal(t, ttl, constants.SessionDuration)
}

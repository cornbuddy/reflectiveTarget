package daos

import (
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/domain/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionStoreShouldReturnNilWhenNoSessionFound(t *testing.T) {
	t.Parallel()

	token := "not-exists"
	got, err := store.Get(ctx, token)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestSessionStoreShouldSaveAndRestoreSession(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc  string
		token string
		want  sessiondata.SessionData
	}

	testCases := []testCase{{
		"should add entry for authenticated user",
		"kek",
		sessiondata.SessionData{
			IsAuthenticated: true,
			Username:        "user",
			UserID:          69,
		},
	}, {
		"should add entry for anonymous user",
		"kek1",
		sessiondata.SessionData{
			IsAuthenticated: false,
		},
	}}

	for _, tc := range testCases {
		err := store.Update(ctx, tc.token, tc.want)
		require.NoError(t, err)

		got, err := store.Get(ctx, tc.token)
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, tc.want, *got)

		key := store.isAuthenticatedKey(tc.token)
		ttl, err := cache.TTL(ctx, key).Result()
		require.NoError(t, err)
		assert.Equal(t, ttl, constants.SessionDuration)
	}
}

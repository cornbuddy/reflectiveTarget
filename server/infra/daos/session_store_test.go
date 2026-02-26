package daos

import (
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/domain/constants"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionStoreShouldReturnNilWhenNoSessionFound(t *testing.T) {
	t.Parallel()

	token := "not-exists"
	got, err := store.IsAuthenitcated(ctx, token)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestSessionStoreShouldSaveAndRestoreSession(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc            string
		token           string
		isAuthenticated bool
	}

	testCases := []testCase{{
		desc:            "should add entry for authenticated user",
		token:           "kek",
		isAuthenticated: true,
	}, {
		desc:            "should add entry for anonymous user",
		token:           "kek1",
		isAuthenticated: false,
	}}

	for _, tc := range testCases {
		err := store.SaveSession(ctx, tc.token, tc.isAuthenticated)
		require.NoError(t, err)

		got, err := store.IsAuthenitcated(ctx, tc.token)
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, tc.isAuthenticated, *got)

		key := store.isAuthenticatedKey(tc.token)
		ttl, err := cache.TTL(ctx, key).Result()
		require.NoError(t, err)
		assert.Equal(t, ttl, constants.SessionDuration)
	}
}

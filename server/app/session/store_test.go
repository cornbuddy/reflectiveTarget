package session_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/session"
)

func TestSessionStoreShouldReturnNilWhenNoSessionFound(t *testing.T) {
	t.Parallel()

	token := session.MakeID()
	got, err := store.Get(ctx, token)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestSessionStoreShouldSaveAndRestoreSession(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc  string
		token session.SessionID
		want  session.Data
	}

	testCases := []testCase{{
		"should add entry for authenticated user",
		session.MakeID(),
		session.Data{userID, username},
	}, {
		"should add entry for anonymous user",
		session.MakeID(),
		session.Data{},
	}}

	for _, tc := range testCases {
		require.NoError(t, store.Update(ctx, tc.token, tc.want))

		got, err := store.Get(ctx, tc.token)
		require.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, tc.want, *got)
	}
}

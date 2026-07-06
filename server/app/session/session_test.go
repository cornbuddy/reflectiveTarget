package session_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

const (
	username = "username"
	userID   = valueobjects.ID(69)
)

func TestContext(t *testing.T) {
	t.Parallel()

}

func TestSessionData(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc string
		data session.Data
		want bool
	}

	testCases := []testCase{{
		"should be false by default",
		session.Data{},
		false,
	}, {
		"should be true if ID is not zero",
		session.Data{UserID: userID},
		true,
	}, {
		"should be true if username is not empty",
		session.Data{Username: username},
		true,
	}, {
		"should be true if both are not zero",
		session.Data{UserID: userID, Username: username},
		true,
	}}

	for _, tc := range testCases {
		got := tc.data.IsAuthenticated()
		assert.Equal(t, tc.want, got, tc.desc)
	}
}

func TestSessionStoreShouldReturnNilWhenNoSessionFound(t *testing.T) {
	t.Parallel()

	token := session.SessionID("not-exists")
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
		"kek",
		session.Data{userID, username},
	}, {
		"should add entry for anonymous user",
		"kek1",
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

func TestMake(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc string
		ctx  context.Context
		want session.Data
	}

	testData := session.Data{userID, username}
	testCases := []testCase{{
		"should return zero value by default",
		ctx,
		session.Data{},
	}, {
		"should return explicitly set logger",
		context.WithValue(ctx, session.SessionDataCtx, &testData),
		testData,
	}}

	for _, tc := range testCases {
		got := session.Read(tc.ctx)
		assert.EqualExportedValues(t, tc.want, got, tc.desc)
	}
}

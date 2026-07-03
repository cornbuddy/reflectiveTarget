package session_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/session"
)

func TestSessionStoreShouldReturnNilWhenNoSessionFound(t *testing.T) {
	t.Parallel()

	token := sessiondata.SessionID("not-exists")
	got, err := store.Get(ctx, token)
	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestSessionStoreShouldSaveAndRestoreSession(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc  string
		token sessiondata.SessionID
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
		want Data
	}

	testData := Data{true, 69, "kek"}
	testCases := []testCase{{
		"should return zero value by default",
		ctx,
		Data{},
	}, {
		"should return explicitly set logger",
		context.WithValue(ctx, SessionDataCtx, &testData),
		testData,
	}}

	for _, tc := range testCases {
		got := Read(tc.ctx)
		assert.EqualExportedValues(t, tc.want, got, tc.desc)
	}
}

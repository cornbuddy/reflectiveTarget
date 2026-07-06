package session_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

const (
	username = "username"
	userID   = valueobjects.ID(69)
)

func TestContext(t *testing.T) {
	t.Parallel()

	want := session.Data{userID, username}
	ctx := session.Context(context.TODO(), &want)
	got := session.Read(ctx)
	assert.EqualExportedValues(t, want, got)
}

func TestRead(t *testing.T) {
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
		"should return session data",
		session.Context(ctx, &testData),
		testData,
	}}

	for _, tc := range testCases {
		got := session.Read(tc.ctx)
		assert.EqualExportedValues(t, tc.want, got, tc.desc)
	}
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

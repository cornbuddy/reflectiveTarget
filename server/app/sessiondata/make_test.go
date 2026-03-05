package sessiondata_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	. "github.com/cornbuddy/reflectiveTarget/server/infra/logger"
)

var (
	ctx = context.TODO()
)

func TestMake(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc string
		ctx  context.Context
		want SessionData
	}

	testLog := Log.Named("test")
	testCases := []testCase{{
		"should return default logger by default",
		ctx,
		SessionData{
			Logger: *Log,
		},
	}, {
		"should return explicitly set logger",
		context.WithValue(ctx, LoggerCtx, *testLog),
		SessionData{
			Logger: *testLog,
		},
	}, {
		"should set authenticated flag",
		context.WithValue(ctx, AuthenticatedCtx, toPtr(true)),
		SessionData{
			IsAuthetnicated: true,
		},
	}, {
		"should set username",
		context.WithValue(ctx, UsernameCtx, toPtr("kek")),
		SessionData{
			Username: "kek",
		},
	}, {
		"should set user id",
		context.WithValue(ctx, UserIDCtx, toPtr(valueobjects.ID(69))),
		SessionData{
			UserID: valueobjects.ID(69),
		},
	}}

	for _, tc := range testCases {
		got := Make(tc.ctx)
		assert.EqualExportedValues(t, tc.want, got, tc.desc)
	}
}

func toPtr[T any](v T) *T { return &v }

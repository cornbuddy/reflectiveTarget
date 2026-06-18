package log_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	want := zap.New(nil)
	ctx := context.WithValue(context.TODO(), log.LoggerCtx, want)
	assert.Same(t, want, log.Logger(ctx), "extracts logger from context")

	got := log.Logger(ctx, zap.String("kek", "kek"))
	assert.NotSame(t, want, got, "decorates logger with fields")

	got = log.Logger(context.TODO())
	assert.NotSame(t, want, got, "falls back to default if context is empty")
}

func TestContext(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc   string
		fields []zap.Field
	}

	testCases := []testCase{{
		"works with empty fields",
		[]zap.Field{},
	}, {
		"works with not empty fields",
		[]zap.Field{zap.String("kek", "kek")},
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx := log.Context(context.TODO(), tc.fields...)
			logger, ok := ctx.Value(log.LoggerCtx).(*zap.Logger)
			assert.True(t, ok)
			assert.NotNil(t, logger)
			assert.IsType(t, &zap.Logger{}, logger)
		})
	}
}

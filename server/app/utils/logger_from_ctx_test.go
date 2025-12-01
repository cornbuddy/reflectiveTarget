package utils

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	. "github.com/cornbuddy/reflectiveTarget/server/infra/logger"
)

func TestShouldReturnLoggerIfContextContainsLogger(t *testing.T) {
	t.Parallel()

	wantLog := Log.With(zap.String("kek", "kek"))
	ctx := context.WithValue(context.TODO(), constants.LoggerCtx, wantLog)
	gotLog := LoggerFromCtx(ctx)
	assert.Equal(t, wantLog, gotLog)
}

func TestShouldReturnDefaultLoggerIfContextIsEmpty(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	log := LoggerFromCtx(ctx)
	assert.NotNil(t, log)
	assert.Equal(t, log, Log)
}

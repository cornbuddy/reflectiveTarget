package utils

import (
	"context"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	. "github.com/cornbuddy/reflectiveTarget/server/infra/logger"
)

func LoggerFromCtx(ctx context.Context) *zap.Logger {
	log := ctx.Value(constants.LoggerCtx)
	if log == nil {
		Log.Warn("missing logger in context", zap.Any("context", ctx))
		return Log
	} else {
		return log.(*zap.Logger)
	}
}

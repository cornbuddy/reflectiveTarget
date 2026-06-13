package log

import (
	"context"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
)

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	Logger(ctx).Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	Logger(ctx).Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	Logger(ctx).Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	Logger(ctx).Error(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	Logger(ctx).Fatal(msg, fields...)
}

// builds logger instance from context
func Logger(ctx context.Context) *zap.Logger {
	reqID, ok := ctx.Value(constants.RequestIDCtx).(string)
	if !ok {
		logger.Warn("missing request ID in context", zap.Any("context", ctx))
		return logger
	} else {
		return logger.With(zap.String("request-id", reqID))
	}
}

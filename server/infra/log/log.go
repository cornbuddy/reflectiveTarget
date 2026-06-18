package log

import (
	"context"

	"go.uber.org/zap"
)

const LoggerCtx = "logger-ctx"

// returns child context with zap logger in it
func Context(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, LoggerCtx, logger.With(fields...))
}

// returns logger from context if defined, else returns default logger
func Logger(ctx context.Context, fields ...zap.Field) *zap.Logger {
	log, ok := ctx.Value(LoggerCtx).(*zap.Logger)
	if !ok {
		logger.Warn(
			"failed to extract logger from context, using default one",
			zap.Any("context", ctx),
		)
		return logger
	}

	return log.With(fields...)
}

package sessiondata

import (
	"context"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	. "github.com/cornbuddy/reflectiveTarget/server/infra/logger"
)

const (
	AuthenticatedCtx = "isAuthenticated"
	LoggerCtx        = "logger"
	UsernameCtx      = "username"
	UserIDCtx        = "userID"
)

type SessionData struct {
	IsAuthetnicated bool
	UserID          valueobjects.ID
	Username        string
	Logger          zap.Logger
}

func Make(ctx context.Context) SessionData {
	return SessionData{
		IsAuthetnicated: getData[bool](ctx, AuthenticatedCtx),
		UserID:          getData[valueobjects.ID](ctx, UserIDCtx),
		Username:        getData[string](ctx, UsernameCtx),
		Logger:          getLogger(ctx),
	}
}

func getLogger(ctx context.Context) zap.Logger {
	log := ctx.Value(LoggerCtx)
	if log == nil {
		Log.Warn("missing logger in context", zap.Any("context", ctx))
		return *Log
	} else {
		return log.(zap.Logger)
	}
}

func getData[T any](ctx context.Context, key string) T {
	var result T
	val := ctx.Value(key)
	if val != nil {
		result = *val.(*T)
	}

	return result
}

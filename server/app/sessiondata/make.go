package sessiondata

import (
	"context"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

const (
	SessionDataCtx = "session-data"
)

type SessionData struct {
	IsAuthetnicated bool
	UserID          valueobjects.ID
	Username        string
}

func Read(ctx context.Context) SessionData {
	return getData[SessionData](ctx, SessionDataCtx)
}

func getData[T any](ctx context.Context, key string) T {
	var result T
	val := ctx.Value(key)
	if val != nil {
		result = *val.(*T)
	}

	return result
}

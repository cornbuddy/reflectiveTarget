package sessiondata

import (
	"context"
	"time"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type ctxKey int

const (
	SessionDataCtx ctxKey = iota
)

const (
	SessionDuration = 30 * 24 * time.Hour
)

type SessionID string

type SessionData struct {
	IsAuthenticated bool
	UserID          valueobjects.ID
	Username        string
}

// reads session data from context. returns zero object if not found
func Read(ctx context.Context) SessionData {
	return getData[SessionData](ctx, SessionDataCtx)
}

func getData[T any](ctx context.Context, key ctxKey) T {
	var result T
	val := ctx.Value(key)
	if val != nil {
		result = *val.(*T)
	}

	return result
}

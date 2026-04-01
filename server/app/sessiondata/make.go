package sessiondata

import (
	"context"
	"time"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

const (
	SessionDataCtx  = "session-data"
	SessionDuration = 30 * 24 * time.Hour
)

type SessionData struct {
	IsAuthenticated bool            `json:"isAuthenticated"`
	UserID          valueobjects.ID `json:"userID"`
	Username        string          `json:"username"`
}

// reads session data from context. returns zero object if not found
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

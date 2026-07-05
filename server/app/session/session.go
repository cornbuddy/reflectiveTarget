package session

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
	// duration of the http session
	Duration = 30 * 24 * time.Hour
)

type SessionID string

type Data struct {
	IsAuthenticated bool
	UserID          valueobjects.ID
	Username        string
}

type Store interface {
	Get(context.Context, SessionID) (*Data, error)
	Update(context.Context, SessionID, Data) error
}

// reads session data from context. returns zero object if not found
func Read(ctx context.Context) Data {
	return getData[Data](ctx, SessionDataCtx)
}

func getData[T any](ctx context.Context, key ctxKey) T {
	var result T
	val := ctx.Value(key)
	if val != nil {
		result = *val.(*T)
	}

	return result
}

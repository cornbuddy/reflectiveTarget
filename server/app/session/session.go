package session

import (
	"context"
	"time"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type ctxKey int

const sessionDataCtx ctxKey = iota

// duration of the http session
const Duration = 30 * 24 * time.Hour

type SessionID string

type Data struct {
	UserID   valueobjects.ID
	Username string
}

type Store interface {
	Get(context.Context, SessionID) (*Data, error)
	Update(context.Context, SessionID, Data) error
}

// returns context with the session data
func Context(ctx context.Context, data *Data) context.Context {
	return context.WithValue(ctx, sessionDataCtx, data)
}

// reads session data from context. returns zero object if not found
func Read(ctx context.Context) Data {
	var result Data
	val := ctx.Value(sessionDataCtx)
	if val != nil {
		result = *val.(*Data)
	}

	return result
}

func (d *Data) IsAuthenticated() bool {
	return d.UserID != 0 || len(d.Username) != 0
}

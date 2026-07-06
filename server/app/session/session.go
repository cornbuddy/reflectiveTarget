package session

import (
	"context"
	"time"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type ctxKey int

const SessionDataCtx ctxKey = iota

const (
	// duration of the http session
	Duration = 30 * 24 * time.Hour
)

type SessionID string

type Data struct {
	UserID   valueobjects.ID
	Username string
}

type Store interface {
	Get(context.Context, SessionID) (*Data, error)
	Update(context.Context, SessionID, Data) error
}

// reads session data from context. returns zero object if not found
func Read(ctx context.Context) Data {
	var result Data
	val := ctx.Value(SessionDataCtx)
	if val != nil {
		result = *val.(*Data)
	}

	return result
}

func (d *Data) IsAuthenticated() bool {
	return d.UserID != 0 || len(d.Username) != 0
}

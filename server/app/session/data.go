package session

import (
	"context"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type ctxKey int

const sessionDataCtx ctxKey = iota

type Data struct {
	UserID   valueobjects.ID
	Username string
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

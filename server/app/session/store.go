package session

import (
	"context"
	"time"
)

// duration of the http session
const Duration = 30 * 24 * time.Hour

type Store interface {
	Get(context.Context, SessionID) (*Data, error)
	Update(context.Context, SessionID, Data) error
}

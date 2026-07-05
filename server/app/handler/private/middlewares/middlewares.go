package middlewares

import (
	"github.com/cornbuddy/reflectiveTarget/server/app/session"
)

type Middleware struct {
	SessionStore session.Store
}

package middlewares

import (
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type Middleware struct {
	SessionStore daos.SessionStore
}

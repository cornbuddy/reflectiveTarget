package middlewares

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

type MiddlewareFunc func(http.Handler) http.Handler

type Middleware struct {
	daos.SessionStore
}

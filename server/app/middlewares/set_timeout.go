package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
)

func (mw Middleware) SetTimeout(timeout time.Duration) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return setTimeout(timeout, next)
	}
}

func setTimeout(timeout time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := utils.LoggerFromCtx(r.Context())
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer func() {
			cancel()
			if ctx.Err() == context.DeadlineExceeded {
				log.Error("request timed out")
				w.WriteHeader(http.StatusGatewayTimeout)
			}
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

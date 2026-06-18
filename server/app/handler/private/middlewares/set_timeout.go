package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

func (mw Middleware) SetTimeout(timeout time.Duration) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return setTimeout(timeout, next)
	}
}

func setTimeout(timeout time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		log := log.Logger(ctx)
		defer func() {
			err := ctx.Err()
			cancel()
			if errors.Is(err, context.DeadlineExceeded) {
				log.Error("request timed out")
				utils.HttpError(w, http.StatusGatewayTimeout)
			}
		}()

		log.Debug("timeout is set", zap.Duration("timeout", timeout))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

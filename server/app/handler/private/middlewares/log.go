package middlewares

import (
	"net/http"
	"time"

	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

// logs http requests. adds logger instance to request's context as a side
// effect
func (mw Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		ctx := log.Context(r.Context(), zap.String("request-id", xid.New().String()))

		next.ServeHTTP(w, r.WithContext(ctx))

		duration := time.Since(startTime)
		log := log.Logger(ctx)
		log.Info("done processing HTTP request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.Duration("duration", duration),
		)
	})
}

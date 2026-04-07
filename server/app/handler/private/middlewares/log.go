package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/xid"
	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	. "github.com/cornbuddy/reflectiveTarget/server/infra/logger"
)

func (mw Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		reqId := xid.New().String()
		log := Log.With(zap.String("request-id", reqId))
		ctx := context.WithValue(r.Context(), constants.LoggerCtx, log)

		next.ServeHTTP(w, r.WithContext(ctx))

		duration := time.Since(startTime)
		log.Info("done processing HTTP request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.Duration("duration", duration),
		)
	})
}

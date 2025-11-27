package middlewares

import (
	"net/http"
	"time"

	. "github.com/cornbuddy/reflectiveTarget/server/app/logger"
	"go.uber.org/zap"
)

func (mw Middleware) Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		elapsedTime := time.Since(startTime)
		Log.Info("HTTP request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.Path),
			zap.Duration("duration", elapsedTime),
		)
	})
}

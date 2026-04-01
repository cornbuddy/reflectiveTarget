package middlewares

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
	"go.uber.org/zap"
)

// checks if user is authenticated, otherwise throws 403
func (mw Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := utils.LoggerFromCtx(ctx)
		session := sessiondata.Read(ctx)

		if session.IsAuthenticated {
			log.Debug("request is authenticated")
			next.ServeHTTP(w, r)
		} else {
			log.Warn("request is unauthenticated", zap.Any("session", session))
			http.Error(w, "forbidden", http.StatusForbidden)
		}
	})
}

package middlewares

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

// checks if user is authenticated, otherwise throws 403
func (mw Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := log.Logger(ctx)
		session := session.Read(ctx)
		if session.IsAuthenticated() {
			log.Debug("request is authenticated")
			next.ServeHTTP(w, r)
		} else {
			log.Warn("request is unauthenticated")
			utils.HttpError(w, http.StatusForbidden)
		}
	})
}

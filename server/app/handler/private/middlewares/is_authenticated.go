package middlewares

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

// checks if user is authenticated, otherwise throws 403
func (mw Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		session := sessiondata.Read(ctx)
		if session.IsAuthenticated {
			log.Debug(ctx, "request is authenticated")
			next.ServeHTTP(w, r)
		} else {
			log.Warn(ctx, "request is unauthenticated")
			http.Error(w, "forbidden", http.StatusForbidden)
		}
	})
}

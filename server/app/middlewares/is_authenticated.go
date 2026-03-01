package middlewares

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
)

// checks if user is authenticated, otherwise throws 403
func (mw Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := utils.LoggerFromCtx(ctx)
		auth, ok := ctx.Value(constants.AuthenticatedCtx).(*bool)

		if !ok {
			log.Warn("request is unauthenticated")
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		if *auth {
			log.Debug("request is authenticated")
			next.ServeHTTP(w, r)
		} else {
			log.Warn("request is unauthenticated")
			http.Error(w, "forbidden", http.StatusForbidden)
		}
	})
}

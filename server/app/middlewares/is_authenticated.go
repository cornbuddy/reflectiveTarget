package middlewares

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	. "github.com/cornbuddy/reflectiveTarget/server/app/logger"
)

func (mw Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(constants.SessionCookieName)
		if err != nil {
			Log.Warn("no session token in request")
			next.ServeHTTP(w, r)
			return
		}

		token := cookie.Value
		value, err := mw.SessionStore.IsAuthenitcated(token)
		if err != nil {
			internalServerError(w, err.Error())
			return
		}

		key := constants.AuthenticatedCtx
		ctx := context.WithValue(r.Context(), key, value)
		Log.Info("context set up",
			zap.String("token", token),
			zap.Bool("is-authenticated", *value),
		)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

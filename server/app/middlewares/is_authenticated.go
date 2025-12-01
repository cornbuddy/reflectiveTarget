package middlewares

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
)

func (mw Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := utils.LoggerFromCtx(ctx)

		cookie, err := r.Cookie(constants.SessionCookieName)
		if err != nil {
			log.Warn("no session token in request")
			next.ServeHTTP(w, r)
			return
		}

		token := cookie.Value
		value, err := mw.SessionStore.IsAuthenitcated(token)
		if err != nil {
			internalServerError(w, err.Error())
			return
		}

		log = log.With(zap.String("token", token))
		if value == nil {
			log.Warn("auth info not included in the context")
		} else {
			log.Info("adding auth info to context",
				zap.Bool("is-authenticated", *value),
			)
			key := constants.AuthenticatedCtx
			ctx = context.WithValue(r.Context(), key, value)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

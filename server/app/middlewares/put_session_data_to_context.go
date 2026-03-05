package middlewares

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
)

// this function reads session data from session store and writes this data to
// request's context
func (mw Middleware) PutSessionDataToContext(next http.Handler) http.Handler {
	panic("todo: put SessionData object into the context")

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
		value, err := mw.SessionStore.IsAuthenticated(ctx, token)
		if err != nil {
			internalServerError(w, err.Error())
			return
		}

		log = log.With(zap.String("token", token))
		if value == nil {
			log.Warn("auth info not included in the context")
		} else {
			log.Debug("adding auth info to context",
				zap.Bool("is-authenticated", *value),
			)
			key := constants.AuthenticatedCtx
			ctx = context.WithValue(ctx, key, value)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

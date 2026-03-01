package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
)

func (mw Middleware) SaveSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := utils.LoggerFromCtx(ctx)
		store := mw.SessionStore
		cookie, err := r.Cookie(constants.SessionCookieName)

		// not empty error means cookie doesn't exist, hence should be
		// set
		if err != nil {
			log.Info("registering new session...")
			_, err := utils.SaveSession(ctx, store, false, w)
			if err != nil {
				internalServerError(w, err.Error())
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		// empty error means cookie exists, hence session token should
		// be validated
		token := cookie.Value
		log = log.With(zap.String("token", token))
		log.Debug("validating session...")
		auth, err := store.IsAuthenticated(ctx, token)
		if err != nil {
			internalServerError(w, err.Error())
			return
		}

		// cookie is present, but not found in the session store.
		// seems like cache key expired earlier than cookie. kinda
		// suspicious, let's reset the session
		if auth == nil {
			log.Warn("session is not registered")
			_, err := utils.SaveSession(ctx, store, false, w)
			if err != nil {
				internalServerError(w, err.Error())
				return
			}
		}

		// session token was either found in the session store, or was
		// set earlier, so let's process the request
		log.Debug("session is validated")
		next.ServeHTTP(w, r)
	})
}

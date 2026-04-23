package middlewares

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
)

// saves session into store. updates request contexts with session data
func (mw Middleware) SaveSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := utils.LoggerFromCtx(ctx)
		store := mw.SessionStore
		cookie, err := r.Cookie(constants.SessionCookieName)

		// not empty error means cookie doesn't exist, hence should be
		// set
		emptySession := sessiondata.SessionData{}
		if err != nil {
			log.Info("registering new session...")
			_, err := utils.SaveSession(ctx, store, emptySession, w)
			if err != nil {
				utils.InternalServerError(log, w, "failed to save session", err)
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		// empty error means cookie exists, hence session sessionId
		// should be validated
		sessionId := cookie.Value
		log = log.With(zap.String("token", sessionId))
		log.Debug("validating session...")
		session, err := store.Get(ctx, sessionId)
		if err != nil {
			utils.InternalServerError(log, w, "failed to fetch session", err)
			return
		}

		// cookie is present, but not found in the session store.
		// seems like cache key expired earlier than cookie. kinda
		// suspicious, let's reset the session
		if session == nil {
			log.Warn("session is not registered")
			session = &emptySession
			_, err := utils.SaveSession(ctx, store, *session, w)
			if err != nil {
				utils.InternalServerError(log, w, "failed to save session", err)
				return
			}
		}

		// session token was either found in the session store, or was
		// set earlier, so let's process the request
		log.Debug("session is validated", zap.Any("session", *session))
		newCtx := context.WithValue(ctx, sessiondata.SessionDataCtx, session)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

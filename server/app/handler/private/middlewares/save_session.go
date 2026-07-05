package middlewares

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

// saves session into store. updates request contexts with session data
func (mw Middleware) SaveSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := log.Logger(ctx)
		store := mw.SessionStore
		cookie, err := r.Cookie(constants.SessionCookieName)

		// not empty error means cookie doesn't exist, hence should be
		// set
		emptySession := session.Data{}
		if err != nil {
			log.Info("registering new session...")
			_, err := utils.SaveSession(ctx, store, emptySession, w)
			if err != nil {
				log.Error("failed to save session", zap.Error(err))
				http.Error(w, "internal server error", http.StatusInternalServerError)

				return
			}

			next.ServeHTTP(w, r)

			return
		}

		// empty error means cookie exists, hence session id
		// should be validated
		id := cookie.Value
		log = log.With(zap.String("token", id))
		log.Debug("validating session...")
		data, err := store.Get(ctx, session.SessionID(id))
		if err != nil {
			log.Error("failed to fetch session", zap.Error(err))
			http.Error(w, "internal server error", http.StatusInternalServerError)

			return
		}

		// cookie is present, but not found in the session store.
		// seems like cache key expired earlier than cookie. kinda
		// suspicious, let's reset the session
		if data == nil {
			log.Warn("session is not registered")
			data = &emptySession
			_, err := utils.SaveSession(ctx, store, *data, w)
			if err != nil {
				log.Error("failed to save session", zap.Error(err))
				http.Error(w, "internal server error", http.StatusInternalServerError)

				return
			}
		}

		// session token was either found in the session store, or was
		// set earlier, so let's process the request
		log.Debug("session is validated", zap.Any("session", *data))
		newCtx := context.WithValue(ctx, session.SessionDataCtx, data)
		next.ServeHTTP(w, r.WithContext(newCtx))
	})
}

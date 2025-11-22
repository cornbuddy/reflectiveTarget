package middlewares

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/cornbuddy/reflectiveTarget/server/app/handlers"
	"github.com/cornbuddy/reflectiveTarget/server/domain/constants"
)

func (mw Middleware) SaveSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// empty error means cookie exists, hence session is active
		if _, err := r.Cookie(handlers.SessionCookieName); err == nil {
			next.ServeHTTP(w, r)
			return
		}

		token := uuid.NewString()
		err := mw.SessionStore.SaveSession(token, false)
		if err != nil {
			status := http.StatusInternalServerError
			http.Error(w, err.Error(), status)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    handlers.SessionCookieName,
			Value:   token,
			Expires: time.Now().Add(constants.SessionDuration),
		})
		next.ServeHTTP(w, r)
	})
}

package middlewares

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/handlers"
	"github.com/cornbuddy/reflectiveTarget/server/app/utils"
)

func (mw Middleware) SaveSession(next http.Handler) http.Handler {
	internalServerError := func(w http.ResponseWriter, msg string) {
		http.Error(w, msg, http.StatusInternalServerError)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(handlers.SessionCookieName)

		// not empty error means cookie doesn't exist, hence should be
		// set
		if err != nil {
			_, err := utils.SaveSession(mw.SessionStore, false, w)
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
		auth, err := mw.SessionStore.IsAuthenitcated(token)
		if err != nil {
			internalServerError(w, err.Error())
			return
		}

		// cookie is present, but not found in the session store.
		// seems like cache key expired earlier than cookie. kinda
		// suspicious, let's reset the session
		if auth == nil {
			_, err := utils.SaveSession(mw.SessionStore, false, w)
			if err != nil {
				internalServerError(w, err.Error())
				return
			}
		}

		// session token was either found in the session store, or was
		// set earlier, so let's process the request
		next.ServeHTTP(w, r)
	})
}

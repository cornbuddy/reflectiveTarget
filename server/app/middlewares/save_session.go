package middlewares

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/cornbuddy/reflectiveTarget/server/app/handlers"
	"github.com/cornbuddy/reflectiveTarget/server/domain/constants"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

func (mw Middleware) SaveSession(next http.Handler) http.Handler {
	error500 := func(w http.ResponseWriter, msg string) {
		http.Error(w, msg, http.StatusInternalServerError)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(handlers.SessionCookieName)

		// not empty error means cookie doesn't exist, hence should be
		// set
		if err != nil {
			err := saveSession(mw.SessionStore, false, w)
			if err != nil {
				error500(w, err.Error())
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
			error500(w, err.Error())
			return
		}

		// cookie is present, but not found in the session store.
		// seems like cache key expired earlier than cookie. kinda
		// suspicious, let's reset the session
		if auth == nil {
			err := saveSession(mw.SessionStore, false, w)
			if err != nil {
				error500(w, err.Error())
				return
			}
		}

		// session token was either found in the session store, or was
		// set earlier, so let's process the request
		next.ServeHTTP(w, r)
	})
}

func saveSession(
	store daos.SessionStore, isAuthenticated bool, w http.ResponseWriter,
) error {

	token := uuid.NewString()
	if err := store.SaveSession(token, isAuthenticated); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    handlers.SessionCookieName,
		Value:   token,
		Expires: time.Now().Add(constants.SessionDuration),
	})

	return nil
}

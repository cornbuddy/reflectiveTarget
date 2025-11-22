package middlewares

import (
	"context"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/handlers"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

const Authenticated = "isAuthenticated"

type Middleware struct {
	daos.SessionStore
}

func (mw Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cookie *http.Cookie
		for _, c := range r.CookiesNamed(handlers.SessionCookieName) {
			if c.Name == handlers.SessionCookieName {
				cookie = c
			}
		}

		if cookie == nil {
			next.ServeHTTP(w, r)
			return
		}

		token := cookie.Value
		value, err := mw.SessionStore.IsAuthenitcated(token)
		if err != nil {
			msg := err.Error()
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), Authenticated, value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

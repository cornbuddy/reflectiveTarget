package middlewares

import (
	"context"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
)

func (mw Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(constants.SessionCookieName)
		if err != nil {
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
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

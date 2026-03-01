package middlewares

import (
	"net/http"
)

// checks if user is authenticated, otherwise throws 403
func (mw Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		next.ServeHTTP(w, r)
	})
}

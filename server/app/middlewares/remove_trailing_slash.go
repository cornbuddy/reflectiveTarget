package middlewares

import (
	"net/http"
)

// removes trailing `/` from the url
func (mw Middleware) RemoveTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (mw Middleware) SetHeader(key, value string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)
			next.ServeHTTP(w, r)
		})
	}
}

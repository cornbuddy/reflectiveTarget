package middlewares

import (
	"net/http"
)

func Chain(mux http.Handler, middlewares ...MiddlewareFunc) http.Handler {
	for _, mw := range middlewares {
		mux = mw(mux)
	}

	return mux
}

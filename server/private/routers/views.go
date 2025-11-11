package routers

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/private/handlers"
)

type ViewsRouter struct {
	handlers.AuthzHandler
	handlers.IndexHandler
}

func (r ViewsRouter) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", r.IndexHandler.Get)
	mux.HandleFunc("GET /login", r.AuthzHandler.GetLogin)
	mux.HandleFunc("GET /signup", r.AuthzHandler.GetSignup)
	mux.HandleFunc("POST /login", r.AuthzHandler.PostLogin)
	mux.HandleFunc("POST /signup", r.AuthzHandler.PostSignup)

	return mux
}

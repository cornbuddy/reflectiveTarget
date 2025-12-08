package handlers

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/app/middlewares"
)

func NewRouter(config *config.Config) http.Handler {
	authz := authzHandler{
		config.UserDao,
		config.SessionStore,
		SignupFormValidator{config.UserDao},
	}
	index := indexHandler{}
	health := healthHandler{config.HealthDao}
	shots := shotsHandler{
		config.ShotsDao,
		ShotsRequestValidator{},
	}
	mw := middlewares.Middleware{
		SessionStore: config.SessionStore,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index.get)
	mux.HandleFunc("GET /login", authz.getLogin)
	mux.HandleFunc("GET /signup", authz.getSignup)
	mux.HandleFunc("GET /logout", authz.getLogout)
	mux.HandleFunc("POST /login", authz.postLogin)
	mux.HandleFunc("POST /signup", authz.postSignup)

	mux.HandleFunc("GET /api/health", health.get)
	mux.HandleFunc("GET /api/target/{targetID}/shots", shots.get)
	mux.HandleFunc("POST /api/target/{targetID}/shots", shots.post)

	return middlewares.Chain(mux,
		mw.IsAuthenticated,
		mw.SaveSession,
		mw.Logger,
	)
}

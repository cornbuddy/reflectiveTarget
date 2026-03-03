package handlers

import (
	"net/http"
	"time"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/app/middlewares"
)

func NewRouter(config *config.Config) http.Handler {
	mux := http.NewServeMux()
	mw := middlewares.Middleware{
		SessionStore: config.SessionStore,
	}

	index := indexHandler{}
	mux.HandleFunc("GET /{$}", index.get)

	authz := authzHandler{
		config.UserDao,
		config.SessionStore,
		SignupFormValidator{config.UserDao},
		LoginFormValidator{config.UserDao},
	}
	mux.HandleFunc("GET /logout", authz.getLogout)
	mux.HandleFunc("GET /login", authz.getLogin)
	mux.HandleFunc("GET /signup", authz.getSignup)
	mux.HandleFunc("POST /login", authz.postLogin)
	mux.HandleFunc("POST /signup", authz.postSignup)

	targets := targetsHandler{}
	targetsMux := http.NewServeMux()
	targetsMux.HandleFunc("GET /", targets.list)
	targetsMux.HandleFunc("POST /", targets.new)
	targetsMux.HandleFunc("PUT /{targetID}", targets.update)
	mux.Handle(
		"/targets/",
		mw.IsAuthenticated(http.StripPrefix("/targets", targetsMux)),
	)

	health := healthHandler{config.HealthDao}
	shots := shotsHandler{
		config.ShotsDao,
		ShotsRequestValidator{},
	}
	mux.HandleFunc("GET /api/health", health.get)
	mux.HandleFunc("GET /api/target/{targetID}/shots", shots.get)
	mux.HandleFunc("POST /api/target/{targetID}/shots", shots.post)

	return middlewares.Chain(mux,
		mw.PutSessionDataToContext,
		mw.SaveSession,
		// TODO: move timeout to configuration block
		mw.SetTimeout(30*time.Second),
		mw.Logger,
	)
}

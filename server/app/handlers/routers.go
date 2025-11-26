package handlers

import (
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/app/middlewares"
	"github.com/cornbuddy/reflectiveTarget/server/domain/validators"
)

func NewRouter(config *config.Config) http.Handler {
	authz := authzHandler{config.UserDao, config.SessionStore}
	index := indexHandler{}
	health := healthHandler{config.HealthDao}
	shots := shotsHandler{
		config.ShotsDao,
		validators.ShotsRequestValidator{},
	}
	mv := middlewares.Middleware{
		SessionStore: config.SessionStore,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index.get)
	mux.HandleFunc("GET /login", authz.getLogin)
	mux.HandleFunc("GET /signup", authz.getSignup)
	mux.HandleFunc("POST /login", authz.postLogin)
	mux.HandleFunc("POST /logout", authz.postLogout)
	mux.HandleFunc("POST /signup", authz.postSignup)
	mux.HandleFunc("GET /api/health", health.get)
	mux.HandleFunc("GET /api/target/{targetID}/shots", shots.get)
	mux.HandleFunc("POST /api/target/{targetID}/shots", shots.post)

	return mv.SaveSession(mux)
}

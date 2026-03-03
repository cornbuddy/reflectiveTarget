package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/app/middlewares"
)

func NewRouter(config *config.Config) http.Handler {
	r := mux.NewRouter()
	mw := middlewares.Middleware{
		SessionStore: config.SessionStore,
	}

	index := indexHandler{}
	r.HandleFunc("/", index.get).Methods(http.MethodGet)

	authz := authzHandler{
		config.UserDao,
		config.SessionStore,
		SignupFormValidator{config.UserDao},
		LoginFormValidator{config.UserDao},
	}
	r.HandleFunc("/logout", authz.getLogout).Methods(http.MethodGet)
	r.HandleFunc("/login", authz.getLogin).Methods(http.MethodGet)
	r.HandleFunc("/login", authz.postLogin).Methods(http.MethodPost)
	r.HandleFunc("/signup", authz.getSignup).Methods(http.MethodGet)
	r.HandleFunc("/signup", authz.postSignup).Methods(http.MethodPost)

	targets := targetsHandler{}
	child := r.PathPrefix("/targets").Subrouter()
	child.Use(mw.IsAuthenticated)
	child.HandleFunc("", targets.list).Methods(http.MethodGet)
	child.HandleFunc("", targets.new).Methods(http.MethodPost)
	child.HandleFunc("/{targetID}", targets.update).Methods(http.MethodPut)

	health := healthHandler{config.HealthDao}
	shots := shotsHandler{
		config.ShotsDao,
		ShotsRequestValidator{},
	}
	r.HandleFunc("/api/health", health.get).Methods(http.MethodGet)
	r.HandleFunc("/api/target/{targetID}/shots", shots.get).
		Methods(http.MethodGet)
	r.HandleFunc("/api/target/{targetID}/shots", shots.post).
		Methods(http.MethodPost)

	r.Use(
		// TODO: move timeout to configuration block
		mw.SetTimeout(30*time.Second),
		mw.Logger,
		mw.SaveSession,
		mw.PutSessionDataToContext,
	)
	// https://stackoverflow.com/a/56937571
	r.NotFoundHandler = r.NewRoute().HandlerFunc(http.NotFound).GetHandler()

	return r
}

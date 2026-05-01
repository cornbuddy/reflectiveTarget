package handler

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/builders"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/middlewares"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
)

var (
	view   = render.View
	layout = render.Layout
)

func MakeHandler(config *config.Config) http.Handler {
	r := mux.NewRouter()
	mw := middlewares.Middleware{
		SessionStore: config.SessionStore,
	}
	r.Use(
		mw.Logger,
		mw.SaveSession,
		// TODO: move timeout to configuration block
		mw.SetTimeout(30*time.Second),
	)

	index := indexHandler{}
	r.HandleFunc("/", index.get).Methods(http.MethodGet)

	authz := authzHandler{
		config.UserDao,
		config.SessionStore,
		validators.SignupFormValidator{UserDao: config.UserDao},
		validators.LoginFormValidator{UserDao: config.UserDao},
	}
	r.HandleFunc("/logout", authz.getLogout).Methods(http.MethodGet)
	r.HandleFunc("/login", authz.getLogin).Methods(http.MethodGet)
	r.HandleFunc("/login", authz.postLogin).Methods(http.MethodPost)
	r.HandleFunc("/signup", authz.getSignup).Methods(http.MethodGet)
	r.HandleFunc("/signup", authz.postSignup).Methods(http.MethodPost)

	targets := targetsHandler{
		repo: config.TargetRepo,
		builder: builders.TargetBuilder{
			Validator: validators.TargetFormValidator{
				TargetRepo: config.TargetRepo,
			},
		},
	}
	child := r.PathPrefix("/targets").Subrouter()
	child.Use(mw.IsAuthenticated)
	child.HandleFunc("", targets.list).Methods(http.MethodGet)
	child.HandleFunc("/new", targets.makeNew).Methods(http.MethodGet)
	child.HandleFunc("/new", targets.saveNew).Methods(http.MethodPost)
	child.HandleFunc("/{id:[0-9]+}", targets.update).Methods(http.MethodPut)

	health := healthHandler{config.HealthDao}
	shots := shotsHandler{
		config.ShotsDao,
		validators.ShotsRequestValidator{},
	}
	r.HandleFunc("/api/health", health.get).Methods(http.MethodGet)
	r.HandleFunc("/api/target/{targetID:[0-9]+}/shots", shots.get).
		Methods(http.MethodGet)
	r.HandleFunc("/api/target/{targetID:[0-9]+}/shots", shots.post).
		Methods(http.MethodPost)

	// https://stackoverflow.com/a/56937571
	r.NotFoundHandler = r.NewRoute().HandlerFunc(http.NotFound).GetHandler()

	return r
}

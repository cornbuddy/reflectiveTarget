package handler

import (
	"net/http"

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
		mw.SetTimeout(config.Timeout),
	)

	get, post, put := http.MethodGet, http.MethodPost, http.MethodPut
	index := indexHandler{}
	r.HandleFunc("/", index.get).Methods(get)

	authz := authzHandler{
		userDao:         config.UserDao,
		sessionStore:    config.SessionStore,
		signupValidator: validators.SignupFormValidator{UserDao: config.UserDao},
		loginValidator:  validators.LoginFormValidator{UserDao: config.UserDao},
	}
	r.HandleFunc("/logout", authz.getLogout).Methods(get)
	r.HandleFunc("/login", authz.getLogin).Methods(get)
	r.HandleFunc("/login", authz.postLogin).Methods(post)
	r.HandleFunc("/signup", authz.getSignup).Methods(get)
	r.HandleFunc("/signup", authz.postSignup).Methods(post)

	targets := targetsHandler{
		repo: config.TargetRepo,
		builder: builders.TargetBuilder{
			Validator: validators.TargetFormValidator{
				TargetRepo: config.TargetRepo,
			},
		},
	}
	t := r.PathPrefix("/targets").Subrouter()
	t.Use(mw.IsAuthenticated)
	t.HandleFunc("", targets.list).Methods(get)
	t.HandleFunc("/new", targets.getNew).Methods(get)
	t.HandleFunc("/new", targets.postNew).Methods(post)
	t.HandleFunc("/{targetID:[0-9]+}", targets.putExisting).Methods(put)
	t.HandleFunc("/{targetID:[0-9]+}", targets.getExisting).Methods(get)

	health := healthHandler{config.HealthDao}
	shots := shotsHandler{
		config.ShotsDao,
		validators.ShotsRequestValidator{},
	}
	api := r.PathPrefix("/api").Subrouter()
	api.Use(mw.SetHeader("Content-Type", "application/json; charset=utf-8"))
	api.HandleFunc("/health", health.get).Methods(get)
	api.HandleFunc("/target/{targetID:[0-9]+}/shots", shots.get).Methods(get)
	api.HandleFunc("/target/{targetID:[0-9]+}/shots", shots.post).Methods(post)

	// https://stackoverflow.com/a/56937571
	r.NotFoundHandler = r.NewRoute().HandlerFunc(http.NotFound).GetHandler()

	return r
}

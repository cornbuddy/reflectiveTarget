package handlers

import (
	"database/sql"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/validators"
)

type ApiRouter struct {
	*sql.DB
	daos.ShotsDao
	validators.ShotsRequestValidator
}

type ViewsRouter struct {
	daos.UserDao
}

func (r ViewsRouter) Routes() http.Handler {
	authz := authzHandler{r.UserDao}
	index := indexHandler{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index.get)
	mux.HandleFunc("GET /login", authz.getLogin)
	mux.HandleFunc("GET /signup", authz.getSignup)
	mux.HandleFunc("POST /login", authz.postLogin)
	mux.HandleFunc("POST /signup", authz.postSignup)

	return mux
}

func (r ApiRouter) Routes() http.Handler {
	health := healthHandler{DB: r.DB}
	shots := shotsHandler{
		ShotsDao:  r.ShotsDao,
		Validator: r.ShotsRequestValidator,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.get)
	mux.HandleFunc("POST /targets/{targetID}/shots", shots.post)

	return mux
}

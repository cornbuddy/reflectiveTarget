package handlers

import (
	"database/sql"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
)

type ApiRouter struct {
	*sql.DB
	daos.UserDao
}

type ViewsRouter struct {
	daos.UserDao
}

func (r ViewsRouter) Routes() http.Handler {
	authz := AuthzHandler{r.UserDao}
	index := IndexHandler{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", index.Get)
	mux.HandleFunc("GET /login", authz.GetLogin)
	mux.HandleFunc("GET /signup", authz.GetSignup)
	mux.HandleFunc("POST /login", authz.PostLogin)
	mux.HandleFunc("POST /signup", authz.PostSignup)

	return mux
}

func (r ApiRouter) Routes() http.Handler {
	health := HealthHandler{DB: r.DB}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", health.Get)

	return mux
}

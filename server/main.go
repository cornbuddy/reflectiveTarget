package main

import (
	"log"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/private/handlers"
)

func main() {
	config, err := MakeConfig()
	if err != nil {
		log.Fatalf("failed to init application: %s", err)
	}

	health := handlers.HealthRouter{
		DB: config.DB,
	}
	authz := handlers.AuthzRouter{
		UserDao: config.UserDao,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /login", authz.GetLogin)
	mux.HandleFunc("GET /signup", authz.GetSignup)
	mux.HandleFunc("POST /login", authz.PostLogin)
	mux.HandleFunc("POST /signup", authz.PostSignup)
	mux.HandleFunc("GET /api/health", health.Get)

	log.Println("starting http server")
	// TODO: provide port as a config
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to start http server: %s", err)
	}
}

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/app/handlers"
)

const timeout = 30 * time.Second

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Println("initializing application...")
	config, err := config.MakeConfig(ctx)
	if err != nil {
		log.Fatalf("failed to init application: %s", err)
	}

	log.Println("starting http server")
	router := handlers.NewRouter(config)
	// TODO: provide port as a config
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("failed to start http server: %s", err)
	}
}

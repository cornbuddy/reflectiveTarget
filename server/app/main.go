package main

import (
	"context"
	"log"
	"net/http"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/app/handlers"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("initializing application...")
	ctx := context.Background()
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

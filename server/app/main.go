package main

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler"
	"github.com/cornbuddy/reflectiveTarget/server/infra/log"
)

func main() {
	ctx := context.Background()
	log := log.Logger(ctx)
	log.Info("initializing application")
	config, err := config.MakeConfig(ctx)
	if err != nil {
		log.Fatal("failed to make configuration", zap.Error(err))
	}

	log.Info("starting http server")
	router := handler.MakeHandler(config)
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("failed to start http server", zap.Error(err))
	}
}

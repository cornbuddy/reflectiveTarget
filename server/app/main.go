package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

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

	handler := handler.MakeHandler(config)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: config.Timeout * 2,
		IdleTimeout:  config.Timeout * 3,
	}
	log.Info("starting http server")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("failed to start http server", zap.Error(err))
	}
}

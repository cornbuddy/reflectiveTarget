package main

import (
	"context"
	"fmt"
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
	cfg, err := config.MakeConfig(ctx)
	if err != nil {
		log.Fatal("failed to make configuration", zap.Error(err))
	}

	handler := handler.MakeHandler(cfg)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: cfg.Timeout * config.WriteTimeoutMultiplier,
		IdleTimeout:  cfg.Timeout * config.IdleTiemoutMultiplier,
	}
	log.Info("starting http server")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("failed to start http server", zap.Error(err))
	}
}

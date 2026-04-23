package main

import (
	"context"
	"net/http"

	"go.uber.org/zap"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler"
	. "github.com/cornbuddy/reflectiveTarget/server/infra/logger"
)

func main() {
	Log.Info("initializing application")
	ctx := context.Background()
	config, err := config.MakeConfig(ctx)
	if err != nil {
		Log.Fatal("failed to make configuration", zap.Error(err))
	}

	Log.Info("starting http server")
	router := handler.MakeHandler(config)
	// TODO: provide port as a config
	if err := http.ListenAndServe(":8080", router); err != nil {
		Log.Fatal("failed to start http server", zap.Error(err))
	}
}

package log

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var logger *zap.Logger

var devConfig = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
	Encoding:         "console",
	EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}

var prodConfig = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
	Development:      false,
	Encoding:         "json",
	EncoderConfig:    zap.NewProductionEncoderConfig(),
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}

func init() {
	environment := "production"
	if env, found := os.LookupEnv("ENVIRONMENT"); found {
		environment = env
	}

	var err error
	if environment == "production" {
		logger, err = prodConfig.Build()
	} else {
		logger, err = devConfig.Build()
	}

	if err != nil {
		log.Fatalf("failed to setup logger: %v", err)
	}

	logger.Info("logger is set up", zap.String("environment", environment))
}

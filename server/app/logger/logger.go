package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var Log *zap.Logger

var devConfig = zap.Config{
	Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
	Development:      true,
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
		Log, err = prodConfig.Build()
	} else {
		Log, err = devConfig.Build()
	}

	if err != nil {
		log.Fatalf("failed to setup logger: %v", err)
	}

	log.Printf("using %s logger\n", environment)
}

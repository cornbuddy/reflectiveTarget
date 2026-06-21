package utils

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	store daos.SessionStore
)

func TestMain(m *testing.M) {
	var code int
	defer os.Exit(code)

	cleanup, cache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Panicf("failed to setup cache: %v", err)
	}

	defer func() {
		if err := cleanup(); err != nil {
			log.Panicf("failed to cleanup cache: %v", err)
		}
	}()

	store = daos.SessionStore{Cache: cache}

	code = m.Run()
}

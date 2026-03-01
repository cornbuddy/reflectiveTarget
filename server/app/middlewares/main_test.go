package middlewares

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
	mw    Middleware
)

func TestMain(m *testing.M) {
	cleanup, cache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup cache: %v", err)
	}

	defer func() {
		if err := cleanup(); err != nil {
			log.Fatalf("failed to cleanup cache: %v", err)
		}
	}()

	store = daos.SessionStore{
		Cache: cache,
	}
	mw = Middleware{
		SessionStore: store,
	}

	os.Exit(m.Run())
}

package middlewares_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/middlewares"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

const (
	username = "username"
)

var (
	ctx       = context.TODO()
	emptyStub = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {},
	)

	store daos.SessionStore
	mw    middlewares.Middleware
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

	store = daos.SessionStore{Cache: cache}
	mw = middlewares.Middleware{SessionStore: store}

	os.Exit(m.Run())
}

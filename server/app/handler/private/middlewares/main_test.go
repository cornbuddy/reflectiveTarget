package middlewares_test

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/middlewares"
	appsession "github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/infra/session"
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

	store appsession.Store
	mw    middlewares.Middleware
)

func TestMain(m *testing.M) {
	cleanup, cache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup cache: %v", err)
	}

	store = session.RedisStore{Cache: cache}
	mw = middlewares.Middleware{SessionStore: store}

	code, err := utils.RunAndCleanup(ctx, m, cleanup)
	if err != nil {
		log.Fatalf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

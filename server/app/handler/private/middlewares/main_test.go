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
	mw = middlewares.Middleware{SessionStore: store}

	code = m.Run()
}

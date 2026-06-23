package utils_test

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
	cleanup, cache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Panicf("failed to setup cache: %v", err)
	}

	store = daos.SessionStore{Cache: cache}

	code, err := utils.RunAndCleanup(ctx, m, cleanup)
	if err != nil {
		log.Panicf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

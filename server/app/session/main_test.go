package session_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	store session.Store
)

func TestMain(m *testing.M) {
	cleanupRedis, redis, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	store = sessionstore.NewRedisStore(redis)

	code, err := utils.RunAndCleanup(ctx, m, cleanupRedis)
	if err != nil {
		log.Fatalf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

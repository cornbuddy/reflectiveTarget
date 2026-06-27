package redis_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/infra/sessionstore/private/redis"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	store redis.RedisStore
)

func TestMain(m *testing.M) {
	cleanupRedis, testRedis, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	store = redis.RedisStore{Cache: testRedis}

	code, err := utils.RunAndCleanup(ctx, m, cleanupRedis)
	if err != nil {
		log.Fatalf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

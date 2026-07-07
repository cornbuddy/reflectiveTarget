package session_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	cache *redis.Client
	db    *sql.DB
)

func TestMain(m *testing.M) {
	cacheCleanup, testCache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	dbCleanup, testDB, err := utils.StartDB(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	db = testDB
	cache = testCache

	code, err := utils.RunAndCleanup(ctx, m, cacheCleanup, dbCleanup)
	if err != nil {
		log.Fatalf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

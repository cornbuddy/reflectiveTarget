package session_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/infra/session"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	redisStore    session.RedisStore
	postgresStore session.PostgresStore
	db            *sql.DB
)

func TestMain(m *testing.M) {
	cleanup, cache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	cleanup, testDB, err := utils.StartDB(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	db = testDB
	redisStore = session.NewPostgresStore{Cache: cache}
	postgresStore = session.NewPostgresStore(db)

	code, err := utils.RunAndCleanup(ctx, m, cleanup)
	if err != nil {
		log.Fatalf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

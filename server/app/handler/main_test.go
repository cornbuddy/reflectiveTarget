package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
	"github.com/cornbuddy/reflectiveTarget/server/infra/utils"
	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.TODO()

	db           *sql.DB
	router       http.HandlerFunc
	sessionStore daos.SessionStore
	userDao      daos.UserDao
)

func TestMain(m *testing.M) {
	cleanupDb, testDb, err := testutils.SetupTestDb(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanupDb(); err != nil {
			log.Fatalf("failed to cleanup db: %v", err)
		}
	}()

	if err := utils.InitDatabase(testDb); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	cleanUpCache, testCache, err := testutils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanUpCache(); err != nil {
			log.Fatalf("failed to clean up cache: %v", err)
		}
	}()

	db = testDb
	router = NewRouter(makeTestConfig(testDb, testCache)).ServeHTTP
	sessionStore = daos.SessionStore{Cache: testCache}
	userDao = daos.UserDao{DB: testDb}

	os.Exit(m.Run())
}

func makeTestConfig(db *sql.DB, cache *redis.Client) *config.Config {
	health := daos.HealthDao{DB: db, Cache: cache}
	shots := daos.ShotsDao{DB: db}
	userDao := daos.UserDao{DB: db}
	sessionStore := daos.SessionStore{Cache: cache}
	targetRepo := repositories.TargetRepo{DB: db}

	return &config.Config{
		HealthDao:    health,
		SessionStore: sessionStore,
		ShotsDao:     shots,
		UserDao:      userDao,
		TargetRepo:   targetRepo,
	}
}

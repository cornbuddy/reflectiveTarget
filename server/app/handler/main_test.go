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
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
	"github.com/redis/go-redis/v9"
)

const (
	ctAppJson = "application/json; charset=utf-8"
	ctForm    = "application/x-www-form-urlencoded"
	get       = http.MethodGet
	post      = http.MethodPost
	put       = http.MethodPut
)

var (
	ctx = context.TODO()

	db           *sql.DB
	router       http.HandlerFunc
	sessionStore daos.SessionStore
	userDao      daos.UserDao
)

func TestMain(m *testing.M) {
	cleanupDb, testDb, err := utils.SetupDB(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanupDb(); err != nil {
			log.Fatalf("failed to cleanup db: %v", err)
		}
	}()

	cleanUpCache, testCache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanUpCache(); err != nil {
			log.Fatalf("failed to clean up cache: %v", err)
		}
	}()

	db = testDb
	router = MakeHandler(makeTestConfig(testDb, testCache)).ServeHTTP
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
		Timeout:      config.DefaultTimeout,
		HealthDao:    health,
		SessionStore: sessionStore,
		ShotsDao:     shots,
		UserDao:      userDao,
		TargetRepo:   targetRepo,
	}
}

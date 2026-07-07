package handler_test

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler"
	appsession "github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
	"github.com/cornbuddy/reflectiveTarget/server/infra/session"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

const (
	ctAppJson = "application/json; charset=utf-8"
	ctForm    = "application/x-www-form-urlencoded"
	get       = http.MethodGet
	post      = http.MethodPost
	put       = http.MethodPut

	username = "username"
	userID   = valueobjects.ID(69)
)

var (
	ctx = context.TODO()

	sessionID   appsession.SessionID
	sessionData appsession.Data
	db          *sql.DB
	router      http.HandlerFunc
	store       appsession.Store
	userDao     daos.UserDao
)

func TestMain(m *testing.M) {
	cleanupDb, testDb, err := utils.SetupDB(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	cleanUpCache, testCache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	db = testDb
	router = handler.MakeHandler(makeTestConfig(testDb, testCache)).ServeHTTP
	store = session.RedisStore{Cache: testCache}
	userDao = daos.UserDao{DB: testDb}

	sessionID = appsession.MakeID()
	sessionData = appsession.Data{UserID: userID, Username: username}
	err = store.Update(ctx, sessionID, sessionData)
	if err != nil {
		log.Fatalf("failed to register session: %v", err)
	}

	code, err := utils.RunAndCleanup(ctx, m, cleanupDb, cleanUpCache)
	if err != nil {
		log.Fatalf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

func makeTestConfig(db *sql.DB, cache *redis.Client) *config.Config {
	health := daos.HealthDao{DB: db, Cache: cache}
	shots := daos.ShotsDao{DB: db}
	userDao := daos.UserDao{DB: db}
	sessionStore := session.RedisStore{Cache: cache}
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

package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/infra/utils"
	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	db           *sql.DB
	userDao      daos.UserDao
	sessionStore daos.SessionStore
	views        http.HandlerFunc
	api          http.HandlerFunc
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
	userDao = daos.UserDao{DB: db}
	sessionStore = daos.SessionStore{
		Ctx:   ctx,
		Cache: testCache,
	}
	views = ViewsRouter{
		UserDao:      userDao,
		SessionStore: sessionStore,
	}.Routes().ServeHTTP
	api = ApiRouter{
		DB:       db,
		ShotsDao: daos.ShotsDao{DB: db},
	}.Routes().ServeHTTP

	os.Exit(m.Run())
}

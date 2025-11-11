package handlers

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/private/validators"
	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
)

var (
	ctx           context.Context
	db            *sql.DB
	userDao       daos.UserDao
	healthHandler HealthHandler
	authzHandler  AuthzHandler
	indexHandler  IndexHandler
	shotsHandler  ShotsHandler
)

func TestMain(m *testing.M) {
	ctx = context.TODO()

	t := &testing.T{}
	cleanup, testDb, err := testdb.SetupTestDb(ctx, t)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	if err := utils.InitDatabase(testDb); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	db = testDb
	userDao = daos.UserDao{DB: db}
	healthHandler = HealthHandler{DB: db}
	authzHandler = AuthzHandler{UserDao: userDao}
	indexHandler = IndexHandler{}
	shotsHandler = ShotsHandler{
		Validator: validators.ShotsRequestValidator{},
		ShotsDao:  daos.ShotsDao{DB: db},
	}

	code := m.Run()
	defer os.Exit(code)

	if err = cleanup(); err != nil {
		log.Fatalf("failed to cleanup test suite: %v", err)
	}

}

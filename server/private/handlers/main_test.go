package handlers

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/utils"
	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
)

var (
	ctx          context.Context
	db           *sql.DB
	userDao      daos.UserDao
	healthRouter HealthRouter
	authzRouter  AuthzRouter
	indexRouter  IndexRouter
	shotsRouter  ShotsRouter
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
	healthRouter = HealthRouter{DB: db}
	authzRouter = AuthzRouter{UserDao: userDao}
	indexRouter = IndexRouter{}
	shotsRouter = ShotsRouter{}

	code := m.Run()
	defer os.Exit(code)

	if err = cleanup(); err != nil {
		log.Fatalf("failed to cleanup test suite: %v", err)
	}

}

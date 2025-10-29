package handlers

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/private/utils"
	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
)

var (
	ctx          context.Context
	db           *sql.DB
	healthRouter HealthRouter
	authzRouter  AuthzRouter
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
	healthRouter = HealthRouter{DB: db}
	authzRouter = AuthzRouter{}

	code := m.Run()
	defer os.Exit(code)

	if err = cleanup(); err != nil {
		log.Fatalf("failed to cleanup test suite: %v", err)
	}

}

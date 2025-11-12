package handlers

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/utils"
	"github.com/cornbuddy/reflectiveTarget/server/private/validators"
	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
)

var (
	ctx     context.Context
	db      *sql.DB
	userDao daos.UserDao
	views   http.HandlerFunc
	api     http.HandlerFunc
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
	views = ViewsRouter{UserDao: userDao}.Routes().ServeHTTP
	api = ApiRouter{
		DB:                    db,
		ShotsDao:              daos.ShotsDao{DB: db},
		ShotsRequestValidator: validators.ShotsRequestValidator{},
	}.Routes().ServeHTTP

	code := m.Run()
	defer os.Exit(code)

	if err = cleanup(); err != nil {
		log.Fatalf("failed to cleanup test suite: %v", err)
	}

}

package validators_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

const (
	freeUsesrname = "username"
	validPassword = "default-password123@"
	weakPassword  = "kekeke"
)

var (
	ctx = context.TODO()

	db      *sql.DB
	userDao daos.UserDao
)

func TestMain(m *testing.M) {
	cleanupDb, testDb, err := utils.SetupDB(ctx)
	if err != nil {
		log.Panicf("failed to setup db: %v", err)
	}

	db = testDb
	userDao = daos.UserDao{DB: testDb}

	code, err := utils.RunAndCleanup(ctx, m, cleanupDb)
	if err != nil {
		log.Panicf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

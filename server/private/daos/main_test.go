package daos

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
	"github.com/cornbuddy/reflectiveTarget/server/private/utils"
	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
)

var (
	ctx      context.Context
	db       *sql.DB
	userDao  UserDao
	password model.Password
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

	pwd, err := model.NewPassword("kek")
	if err != nil {
		log.Fatalf("failed to create password: %v", err)
	}

	db = testDb
	userDao = UserDao{DB: db}
	password = *pwd

	code := m.Run()
	defer os.Exit(code)

	if err = cleanup(); err != nil {
		log.Fatalf("failed to cleanup test suite: %v", err)
	}

}

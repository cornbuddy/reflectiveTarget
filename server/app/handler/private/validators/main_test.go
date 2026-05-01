package validators

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	db      *sql.DB
	userDao daos.UserDao
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

	db = testDb
	userDao = daos.UserDao{DB: testDb}

	os.Exit(m.Run())
}

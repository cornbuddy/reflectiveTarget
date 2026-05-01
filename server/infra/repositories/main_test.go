package repositories_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	db *sql.DB
)

func TestMain(m *testing.M) {
	cleanup, testDb, err := utils.SetupDB(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanup(); err != nil {
			log.Fatalf("failed to clean up db: %v", err)
		}
	}()

	db = testDb

	os.Exit(m.Run())
}

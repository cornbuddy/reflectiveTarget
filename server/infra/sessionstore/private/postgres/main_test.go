package postgres_test

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
	cleanup, testDB, err := utils.StartDB(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	db = testDB

	code, err := utils.RunAndCleanup(ctx, m, cleanup)
	if err != nil {
		log.Fatalf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

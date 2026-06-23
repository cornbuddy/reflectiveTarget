package repositories_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
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

	db = testDb

	code, err := utils.RunAndCleanup(ctx, m, cleanup)
	if err != nil {
		log.Fatalf("failed to cleanup: %v", err)
	}

	os.Exit(code)
}

func makeTarget(owner entities.User) aggregations.Target {
	return aggregations.Target{
		Name:  utils.MakeRandomString(5),
		Owner: owner,
		Questions: valueobjects.Questions{{
			Text: utils.MakeRandomString(10),
		}, {
			Text: utils.MakeRandomString(10),
		}},
		Shots: valueobjects.Shots{},
	}
}

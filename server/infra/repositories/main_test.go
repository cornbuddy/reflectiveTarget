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

	defer func() {
		if err := cleanup(); err != nil {
			log.Fatalf("failed to clean up db: %v", err)
		}
	}()

	db = testDb

	os.Exit(m.Run())
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

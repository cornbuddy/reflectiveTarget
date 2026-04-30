package builders_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	targetRepo repositories.TargetRepo
)

func TestMain(m *testing.M) {
	cleanup, db, err := utils.SetupDB(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanup(); err != nil {
			log.Fatalf("failed to cleanup db: %v", err)
		}
	}()

	targetRepo = repositories.TargetRepo{DB: db}

	os.Exit(m.Run())
}

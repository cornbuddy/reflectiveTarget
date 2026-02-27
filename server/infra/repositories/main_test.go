package repositories_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
	"github.com/cornbuddy/reflectiveTarget/server/infra/utils"
	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	targetRepo repositories.TargetRepo
)

func TestMain(m *testing.M) {
	cleanUpDb, testDb, err := testutils.SetupTestDb(ctx)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanUpDb(); err != nil {
			log.Fatalf("failed to clean up db: %v", err)
		}
	}()

	if err := utils.InitDatabase(testDb); err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	targetRepo = repositories.TargetRepo{testDb}

	os.Exit(m.Run())
}

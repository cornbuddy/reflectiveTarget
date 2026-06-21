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
		log.Panicf("failed to setup db: %v", err)
	}

	defer func() {
		if err := cleanup(); err != nil {
			log.Panicf("failed to cleanup db: %v", err)
		}
	}()

	targetRepo = repositories.TargetRepo{DB: db}

	code := m.Run()
	if code != 0 {
		os.Exit(code)
	}
}

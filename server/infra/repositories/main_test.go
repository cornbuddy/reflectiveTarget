package repositories_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
	"github.com/cornbuddy/reflectiveTarget/server/infra/utils"
	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx = context.TODO()

	db *sql.DB
)

type TargetRepoTests struct {
	targets aggregations.Targets
	repo    *repositories.TargetRepo
	// owns 1 targets
	owner1 *entities.User
	// owns 2 targets
	owner2 *entities.User
}

func (s *TargetRepoTests) PreGroup(t *testgroup.T) {
	owner1, err := entities.NewUser("user1", "password")
	t.Require.NoError(err)

	owner2, err := entities.NewUser("user2", "password")
	t.Require.NoError(err)

	users := entities.Users{*owner1, *owner2}
	t.Require.NoError(testutils.InsertUsers(db, users))

	questions := vo.Questions{
		{Text: "kek1?"},
		{Text: "kek2?"},
	}
	shots := vo.Shots{
		{X: 1, Y: 100},
		{X: 100, Y: 1},
	}
	targets := aggregations.Targets{aggregations.Target{
		Name:      "test1",
		Owner:     users[0],
		Questions: append(vo.Questions{}, questions...),
		Shots:     append(vo.Shots{}, shots...),
	}, {
		Name:      "test2",
		Owner:     users[1],
		Questions: append(vo.Questions{}, questions...),
		Shots:     append(vo.Shots{}, shots...),
	}, {
		Name:      "test3",
		Owner:     users[1],
		Questions: append(vo.Questions{}, questions...),
		Shots:     append(vo.Shots{}, shots...),
	}}
	t.Require.NoError(testutils.InsertTargets(db, targets))

	s.owner1 = &users[0]
	s.owner2 = &users[1]
	s.repo = &repositories.TargetRepo{db}
	s.targets = targets
}

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

	db = testDb

	os.Exit(m.Run())
}

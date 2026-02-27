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

	db         *sql.DB
	targetRepo repositories.TargetRepo
)

type TargetGetTests struct {
	*aggregations.Target
}

func (suite *TargetGetTests) PreGroup(t *testgroup.T) {
	owner, err := entities.NewUser("username", "password")
	t.Require.NoError(err)
	t.Require.NoError(testutils.InsertUser(db, owner))

	suite.Target = &aggregations.Target{
		Name:  "test",
		Owner: *owner,
	}
	t.Require.NoError(testutils.InsertTarget(db, suite.Target))

	targetId := suite.Target.ID
	questions := vo.Questions{
		{Text: "kek1?"},
		{Text: "kek2?"},
	}
	t.Require.NoError(testutils.InsertQuestions(db, &questions, targetId))

	shots := vo.Shots{
		{X: 1, Y: 100},
		{X: 100, Y: 1},
	}
	t.Require.NoError(testutils.InsertShots(db, &shots, targetId, "kek"))

	suite.Target.Questions = questions
	suite.Target.Shots = shots
}

type TargetSaveTests struct{}

type TargetListTargetNames struct{}

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
	targetRepo = repositories.TargetRepo{testDb}

	os.Exit(m.Run())
}

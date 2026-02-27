package repositories_test

import (
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
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

func (*TargetGetTests) ShouldReturnNilIfTargetDoesNotExist(t *testgroup.T) {
	got, err := targetRepo.Get(ctx, 69)
	t.Require.NoError(err)
	t.Nil(got)
}

func (suite *TargetGetTests) ShouldReturnTargetIfExist(t *testgroup.T) {
	got, err := targetRepo.Get(ctx, suite.Target.ID)
	t.Require.NoError(err)
	t.NotNil(got)
	t.EqualValues(*suite.Target, *got)
}

type TargetSaveTests struct{}

func (*TargetSaveTests) ShouldSaveTarget(t *testgroup.T) {
	t.Require.NoError(targetRepo.Save(ctx, nil))
}

type TargetListTargetNames struct{}

func (*TargetListTargetNames) ShouldListTargetNames(t *testgroup.T) {
	targets, err := targetRepo.ListTargetNamesOfUser(ctx, "kek")
	t.Require.NoError(err)
	t.Empty(targets)
}

func TestTargetRepoListTargetNamesForUser(t *testing.T) {
	t.Parallel()
	t.Skip()

	testgroup.RunInParallel(t, new(TargetListTargetNames))
}

func TestTargetRepoGetTargetById(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetGetTests))
}

func TestTargetRepoSaveTarget(t *testing.T) {
	t.Parallel()
	t.Skip()

	testgroup.RunInParallel(t, new(TargetSaveTests))
}

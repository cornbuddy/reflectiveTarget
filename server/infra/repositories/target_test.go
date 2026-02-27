package repositories_test

import (
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type TargetGetTests struct {
	*aggregations.Target
}

func (s *TargetGetTests) PreGroup(t *testgroup.T) {
	owner, err := entities.NewUser("username", "password")
	t.Require.NoError(err)
	t.Require.NoError(testutils.InsertUser(db, owner))

	s.Target = &aggregations.Target{
		Name:  "test",
		Owner: *owner,
	}
	t.Require.NoError(testutils.InsertTarget(db, s.Target))

	question := &valueobjects.Question{Text: "kek?"}
	t.Require.NoError(testutils.InsertQuestion(db, question, s.Target.ID))

	s.Target.Questions = valueobjects.Questions{*question}
}

func (*TargetGetTests) ShouldReturnNilIfTargetDoesNotExist(t *testgroup.T) {
	got, err := targetRepo.Get(ctx, 69)
	t.Require.NoError(err)
	t.Nil(got)
}

func (s *TargetGetTests) ShouldReturnTargetIfExist(t *testgroup.T) {
	got, err := targetRepo.Get(ctx, s.Target.ID)
	t.Require.NoError(err)
	t.NotNil(got)
	t.EqualValues(*s.Target, *got)
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

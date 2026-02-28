package repositories_test

import (
	"testing"

	"github.com/bloomberg/go-testgroup"
)

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

func (*TargetSaveTests) ShouldSaveTarget(t *testgroup.T) {
	t.Require.NoError(targetRepo.Save(ctx, nil))
}

func (*TargetListTargetNames) ShouldListTargetNames(t *testgroup.T) {
	targets, err := targetRepo.ListTargetNamesOfUser(ctx, "kek")
	t.Require.NoError(err)
	t.Empty(targets)
}

func TestTargetRepoListTargetNamesForUser(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetListTargetNames))
}

func TestTargetRepoSaveTarget(t *testing.T) {
	t.Parallel()
	t.Skip()

	testgroup.RunInParallel(t, new(TargetSaveTests))
}

func TestTargetRepoGetTargetById(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetGetTests))
}

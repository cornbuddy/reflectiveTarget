package repositories_test

import (
	"testing"

	"github.com/bloomberg/go-testgroup"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
)

type TargetGetTests struct {
	*aggregations.Target
}

func (s *TargetGetTests) PreGroup(t *testgroup.T) {
	s.Target = &aggregations.Target{}
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

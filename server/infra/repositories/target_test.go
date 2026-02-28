package repositories_test

import (
	"testing"

	"github.com/bloomberg/go-testgroup"
)

func (s *TargetRepoTests) ShouldReturnNilIfTargetDoesNotExist(t *testgroup.T) {
	got, err := s.repo.Get(ctx, 69)
	t.Require.NoError(err)
	t.Nil(got)
}

func (s *TargetRepoTests) ShouldReturnTargetIfExist(t *testgroup.T) {
	target := (*s.targets)[0]
	got, err := s.repo.Get(ctx, target.ID)
	t.Require.NoError(err)
	t.NotNil(got)
	t.EqualValues(target, *got)
}

func (s *TargetRepoTests) ShouldReturnEmptyListIfNoUser(t *testgroup.T) {
	targets, err := s.repo.ListTargetNamesOfUser(ctx, "kek")
	t.Require.NoError(err)
	t.Empty(targets)
}

func (s *TargetRepoTests) ShouldReturnNamesIfTargetsExist(t *testgroup.T) {
	targets, err := s.repo.ListTargetNamesOfUser(ctx, "kek")
	t.Require.NoError(err)
	t.Empty(targets)
}

func (*TargetRepoTests) ShouldSaveTarget(t *testgroup.T) {
	t.Skip("not implemented")
}

func TestTargetRepo(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetRepoTests))
}

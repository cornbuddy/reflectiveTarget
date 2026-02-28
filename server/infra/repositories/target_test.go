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
	want := s.targets[0]
	got, err := s.repo.Get(ctx, want.ID)
	t.Require.NoError(err)
	t.NotNil(got)
	t.EqualValues(want, *got)
}

func (s *TargetRepoTests) ShouldReturnEmptyListIfNoUser(t *testgroup.T) {
	targets, err := s.repo.ListTargetNamesOfUser(ctx, "kek")
	t.Require.NoError(err)
	t.Empty(targets)
}

func (s *TargetRepoTests) ShouldReturnNamesIfTargetsExist(t *testgroup.T) {
	got, err := s.repo.ListTargetNamesOfUser(ctx, s.owner1.Username)
	t.Require.NoError(err)
	t.Len(got, 1)

	got, err = s.repo.ListTargetNamesOfUser(ctx, s.owner2.Username)
	t.Require.NoError(err)
	t.Len(got, 2)
}

func (*TargetRepoTests) ShouldSaveTarget(t *testgroup.T) {
	t.Skip("not implemented")
}

func TestTargetRepo(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetRepoTests))
}

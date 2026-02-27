package repositories_test

import (
	"testing"

	"github.com/bloomberg/go-testgroup"
)

type TargetGetTests struct{}

func (s TargetGetTests) ShouldGetByIdIfExists(t *testgroup.T) {
	_, err := targetRepo.Get(ctx, 69)
	t.Require.NoError(err)
}

type TargetSaveTests struct{}

func (s TargetSaveTests) ShouldSaveTarget(t *testgroup.T) {
	t.Require.NoError(targetRepo.Save(ctx, nil))
}

type TargetListTargetNames struct{}

func (s TargetListTargetNames) ShouldListTargetNames(t *testgroup.T) {
	targets, err := targetRepo.ListTargetNamesOfUser(ctx, "kek")
	t.Require.NoError(err)
	t.Empty(targets)
}

func TestTargetRepoListTargetNamesForUser(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetListTargetNames))
}

func TestTargetRepoGetTargetById(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetGetTests))
}

func TestTargetRepoSaveTarget(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetSaveTests))
}

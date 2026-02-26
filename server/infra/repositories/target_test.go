package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTargetRepoShouldSaveTarget(t *testing.T) {
	t.Parallel()

	require.NoError(t, target.Save(ctx, nil))
}

func TestTargetRepoShouldListTargetNamesForUser(t *testing.T) {
	t.Parallel()

	_, err := target.ListTargetNamesOfUser(ctx, "kek")
	require.NoError(t, err)
}

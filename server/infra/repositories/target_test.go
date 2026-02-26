package repositories_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	// consider using the module below; each repo method belongs to its
	// own suite
	// "github.com/stretchr/testify/suite"
)

func TestTargetRepoShouldGetTargetById(t *testing.T) {
	t.Parallel()

	_, err := target.Get(ctx, 69)
	require.NoError(t, err)
}

func TestTargetRepoShouldSaveTarget(t *testing.T) {
	t.Parallel()

	require.NoError(t, target.Save(ctx, nil))
}

func TestTargetRepoShouldListTargetNamesForUser(t *testing.T) {
	t.Parallel()

	_, err := target.ListTargetNamesOfUser(ctx, "kek")
	require.NoError(t, err)
}

package utils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/infra/utils"
	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestInitDatabase(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	cleanup, db, err := testutils.StartDB(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, cleanup())
	})

	require.NoError(t, utils.InitDatabase(ctx, db))

	tables := []string{"users", "targets", "questions", "shots"}
	for _, table := range tables {
		name := "check if table exists: " + table
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			//nolint:gosec
			query := "SELECT id FROM " + table
			rows, err := db.QueryContext(ctx, query)
			require.NoError(t, err)
			require.NoError(t, rows.Err())
		})
	}
}

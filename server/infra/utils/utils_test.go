package utils

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestInitDatabase(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	cleanup, db, err := utils.SetupTestDb(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, cleanup())
	})

	require.NoError(t, InitDatabase(db))

	tables := []string{"users", "targets", "questions", "shots"}
	for _, table := range tables {
		query := fmt.Sprintf("select * from %s", table)
		_, err := db.Query(query)
		assert.NoError(t, err)
	}
}

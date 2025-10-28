package utils

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"

	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
)

func TestInitDatabase(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	dbContainer, err := testdb.Setup(ctx)
	assert.NoError(t, err)

	t.Cleanup(func() {
		err := testcontainers.TerminateContainer(dbContainer)
		assert.NoError(t, err)
	})

	connStr, err := dbContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	db, err := sql.Open("pgx", connStr)
	assert.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		assert.NoError(t, err)
	})

	err = InitDatabase(db)
	assert.NoError(t, err)

	tables := []string{"users", "targets", "questions", "shots"}
	for _, table := range tables {
		query := fmt.Sprintf("select * from %s", table)
		_, err := db.Query(query)
		assert.NoError(t, err)
	}
}

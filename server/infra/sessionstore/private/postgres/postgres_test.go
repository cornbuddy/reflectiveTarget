package postgres_test

import (
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/infra/sessionstore/private/postgres"
)

const (
	schema = "public"
	table  = "session"
)

type PostgresStoreTest struct {
	store *postgres.PostgresStore
}

func (s *PostgresStoreTest) CreatesProperTable(t *testgroup.T) {
	var exists bool
	q := `SELECT EXISTS (
		SELECT FROM pg_tables
		WHERE schemaname = $1 AND tablename = $2
	)`
	err := db.QueryRowContext(ctx, q, schema, table).Scan(&exists)
	t.Require.NoError(err)
	t.True(exists)
}

func (s *PostgresStoreTest) PreGroup(t *testgroup.T) {
	store, err := postgres.NewPostgresStore(ctx, db)
	t.Require.NoError(err)

	s.store = store
}

func TestNewPostgresStore(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(PostgresStoreTest))
}

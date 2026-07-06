package session_test

import (
	"database/sql"
	"testing"

	"github.com/bloomberg/go-testgroup"

	appsession "github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/infra/session"
)

const (
	schema = "public"
	table  = "sessions"
)

type PostgresStoreTest struct {
	db    *sql.DB
	store *session.PostgresStore
}

func (s *PostgresStoreTest) UpdateCreatesRecodrs(t *testgroup.T) {
	id := appsession.SessionID("kekeke")
	data := appsession.Data{}
	t.Require.NoError(s.store.Update(ctx, id, data))
}

func (s *PostgresStoreTest) TableHasProperColumns(t *testgroup.T) {
	type testCase struct {
		column   string
		dataType string
		nullable string
		defaults *string
	}

	const no, yes = "NO", "YES"
	testCases := []testCase{{
		"id",
		"uuid",
		no,
		new("gen_random_uuid()"),
	}, {
		"session_data",
		"jsonb",
		no,
		new("'{}'::jsonb"),
	}, {
		"ip_address",
		"inet",
		no,
		nil,
	}, {
		"user_id",
		"integer",
		yes,
		nil,
	}, {
		"username",
		"character varying",
		yes,
		nil,
	}}

	const query = `
	SELECT data_type, is_nullable, column_default
	FROM information_schema.columns
	WHERE table_schema = $1 AND table_name = $2 AND column_name = $3`
	for _, tc := range testCases {
		t.Run(tc.column, func(t *testgroup.T) {
			t.Parallel()

			var dataType string
			var nullable string
			var defaults *string
			err := db.QueryRowContext(ctx, query, schema, table, tc.column).
				Scan(&dataType, &nullable, &defaults)
			t.Require.NoError(err)
			t.Equal(tc.dataType, dataType)
			t.Equal(tc.nullable, nullable)
			t.Equal(tc.defaults, defaults)
		})
	}
}

func (s *PostgresStoreTest) CreatesSessionTable(t *testgroup.T) {
	var exists bool
	query := `
	SELECT EXISTS (
		SELECT FROM pg_tables
		WHERE schemaname = $1 AND tablename = $2
	)`
	err := db.QueryRowContext(ctx, query, schema, table).Scan(&exists)
	t.Require.NoError(err)
	t.True(exists)
}

func (s *PostgresStoreTest) PreGroup(t *testgroup.T) {
	store, err := session.NewPostgresStore(ctx, db)
	t.Require.NoError(err)

	s.db = db
	s.store = store
}

func TestNewPostgresStore(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(PostgresStoreTest))
}

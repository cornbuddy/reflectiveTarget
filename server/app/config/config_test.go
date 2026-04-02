package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// those tests are not prarallel because os.Setenv sets env var globally, across
// all goroutines, which messess up test cases when I don't expect env vars to
// be set

func TestInitShouldReturnConfigWhenEnvVarsAreSet(t *testing.T) {
	setEnvVars(t)

	config, err := MakeConfig(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, config)

	val := reflect.ValueOf(*config)
	typ := reflect.TypeFor[Config]()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath == "" {
			v := val.Field(i).Interface()
			assert.NotEmpty(t, v, "field %s should not be empty", field.Name)
		}
	}

	db := config.HealthDao.DB
	cache := config.HealthDao.Cache

	t.Cleanup(func() {
		assert.NoError(t, db.Close())
		assert.NoError(t, cache.Close())
	})

	assert.NoError(t, db.Ping())
	assert.NoError(t, cache.Ping(ctx).Err())

	tables := []string{"users", "targets", "questions", "shots"}
	for _, table := range tables {
		query := fmt.Sprintf("SELECT * FROM %s", table)
		rows, err := db.Query(query)
		require.NoError(t, err)

		cols, err := rows.Columns()
		assert.NoError(t, err)
		assert.NotEmpty(t, cols)
	}
}

func TestInitShouldReturnErrorWhenRequiredEnvVarsAreNotSet(t *testing.T) {
	_, err := MakeConfig(ctx)
	require.ErrorContains(t, err, "is not set")
}

func setEnvVars(t *testing.T) {
	type environmentVariable struct {
		key   string
		value string
	}

	envVars := []environmentVariable{
		{"CACHE_ADDRESS", cacheAddr},
		{"PGPASSWORD", password},
		{"PGUSER", username},
		{"PGDATABASE", database},
		{"PGHOST", dbHost},
	}

	for _, envVar := range envVars {
		t.Setenv(envVar.key, envVar.value)
	}
}

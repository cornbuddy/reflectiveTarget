package config_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/config"
)

// those tests are not prarallel because os.Setenv sets env var globally, across
// all goroutines, which messess up test cases when I don't expect env vars to
// be set

func TestInitSetsTimeout(t *testing.T) {
	setRequiredEnvVars(t)

	type testCase struct {
		desc        string
		envTimeout  *string
		wantTimeout time.Duration
	}

	testCases := []testCase{{
		"should have proper defaults",
		nil,
		config.DefaultTimeout,
	}, {
		"should respect if set",
		new("30"),
		30 * time.Second,
	}, {
		"should ignore crap",
		new("kek"),
		config.DefaultTimeout,
	}, {
		"should switch to default if 0",
		new("0"),
		config.DefaultTimeout,
	}, {
		"should switch to default value is negative",
		new("-69"),
		config.DefaultTimeout,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			defer require.NoError(t, os.Unsetenv("TIMEOUT_SECONDS"))

			if tc.envTimeout != nil {
				t.Setenv("TIMEOUT_SECONDS", *tc.envTimeout)
			}

			got, err := config.MakeConfig(ctx)
			require.NoError(t, err)
			assert.Equal(t, tc.wantTimeout, got.Timeout)
		})
	}
}

func TestInitShouldReturnConfigWhenEnvVarsAreSet(t *testing.T) {
	setRequiredEnvVars(t)

	cfg, err := config.MakeConfig(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, cfg)

	val := reflect.ValueOf(*cfg)
	typ := reflect.TypeFor[config.Config]()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath == "" {
			v := val.Field(i).Interface()
			assert.NotEmpty(t, v, "field %s should not be empty", field.Name)
		}
	}

	db := cfg.HealthDao.DB
	cache := cfg.HealthDao.Cache

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
	_, err := config.MakeConfig(ctx)
	require.ErrorContains(t, err, "is not set")
}

func setRequiredEnvVars(t *testing.T) {
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

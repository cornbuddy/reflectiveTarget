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

func TestOptionalConfiguration(t *testing.T) {
	setRequiredEnvVars(t)

	type extractor func(config.Config) any
	type testCase struct {
		extractor

		desc     string
		envKey   string
		envValue *string
		want     any
	}

	const port = "PORT"
	const timeout = "TIMEOUT_SECONDS"

	portExtr := func(c config.Config) any { return c.Port }
	tiemoutExtr := func(c config.Config) any { return c.Timeout }
	testCases := []testCase{{
		portExtr,
		"should have proper defaults",
		port,
		nil,
		config.DefaultPort,
	}, {
		portExtr,
		"should ignore crap",
		port,
		new("crap"),
		config.DefaultPort,
	}, {
		portExtr,
		"should fallback to default if < 1024",
		port,
		new("80"),
		config.DefaultPort,
	}, {
		portExtr,
		"should fallback to default if > 65535",
		port,
		new("65536"),
		config.DefaultPort,
	}, {
		portExtr,
		"should respect if set",
		port,
		new("42069"),
		42069,
	}, {
		tiemoutExtr,
		"should have proper defaults",
		timeout,
		nil,
		config.DefaultTimeout,
	}, {
		tiemoutExtr,
		"should respect if set",
		timeout,
		new("30"),
		30 * time.Second,
	}, {
		tiemoutExtr,
		"should ignore crap",
		timeout,
		new("kek"),
		config.DefaultTimeout,
	}, {
		tiemoutExtr,
		"should switch to default if 0",
		timeout,
		new("0"),
		config.DefaultTimeout,
	}, {
		tiemoutExtr,
		"should switch to default value is negative",
		timeout,
		new("-69"),
		config.DefaultTimeout,
	}}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("[%s] %s", tc.envKey, tc.desc), func(t *testing.T) {
			defer require.NoError(t, os.Unsetenv(tc.envKey))

			if tc.envValue != nil {
				t.Setenv(tc.envKey, *tc.envValue)
			}

			got, err := config.MakeConfig(ctx)
			require.NoError(t, err)
			assert.Equal(t, tc.want, tc.extractor(*got))
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
	for i := range val.NumField() {
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
		// I don't care about sql injections in tests
		//nolint:gosec
		query := "SELECT * FROM " + table
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

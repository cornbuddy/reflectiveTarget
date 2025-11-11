package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

// those tests are not prarallel because os.Setenv sets env var globally, across
// all goroutines, which messess up test cases when I don't expect env vars to
// be set

func TestMakeHttpHandler(t *testing.T) {
	type testCase struct {
		url        string
		method     string
		statusCode int
	}

	cleanup, err := setDbEnvVars()
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, cleanup())
	})

	config, err := MakeConfig()
	require.NoError(t, err)

	handle := MakeMux(config).ServeHTTP

	testCases := []testCase{{
		url:        "/",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/kek",
		method:     http.MethodGet,
		statusCode: http.StatusNotFound,
	}, {
		url:        "/login",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/login",
		method:     http.MethodPost,
		statusCode: http.StatusBadRequest,
	}, {
		url:        "/signup",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/signup",
		method:     http.MethodPost,
		statusCode: http.StatusBadRequest,
	}, {
		url:        "/api/health",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/api/target/1/shots",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/api/target/1/shots",
		method:     http.MethodPost,
		statusCode: http.StatusBadRequest,
	}}

	for _, tc := range testCases {
		resp := utils.MakeRequest("", tc.method, tc.url, handle, nil)
		msg := fmt.Sprintf("%s %s", tc.method, tc.url)
		assert.Equal(t, tc.statusCode, resp.StatusCode, msg)
	}
}

func TestInitShouldReturnConfigWhenEnvVarsAreSet(t *testing.T) {
	cleanup, err := setDbEnvVars()
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, cleanup())
	})

	config, err := MakeConfig()
	require.NoError(t, err)
	require.NotEmpty(t, config)

	t.Cleanup(func() {
		assert.NoError(t, config.DB.Close())
	})

	assert.NoError(t, config.DB.Ping())

	tables := []string{"users", "targets", "questions", "shots"}
	for _, table := range tables {
		query := fmt.Sprintf("SELECT * FROM %s", table)
		rows, err := config.DB.Query(query)
		require.NoError(t, err)

		cols, err := rows.Columns()
		assert.NoError(t, err)
		assert.NotEmpty(t, cols)
	}
}

func TestInitShouldReturnErrorWhenRequiredEnvVarsAreNotSet(t *testing.T) {
	_, err := MakeConfig()
	require.ErrorIs(t, err, ErrNoEnvVar)
}

type cleanup func() error

type environmentVariable struct {
	key   string
	value string
}

func setDbEnvVars() (cleanup, error) {
	emptyCleanup := func() error { return nil }

	envVars := []environmentVariable{
		{"PGPASSWORD", password},
		{"PGUSER", username},
		{"PGDATABASE", database},
		{"PGHOST", host},
	}

	for _, envVar := range envVars {
		if err := os.Setenv(envVar.key, envVar.value); err != nil {
			return emptyCleanup, err
		}
	}

	clean := func() error {
		for _, envVar := range envVars {
			if err := os.Unsetenv(envVar.key); err != nil {
				return err
			}
		}

		return nil
	}

	return clean, nil
}

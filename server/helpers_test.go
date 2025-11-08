package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/cornbuddy/reflectiveTarget/server/test/db"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestMakeHttpHandler(t *testing.T) {
	t.Parallel()

	type testCase struct {
		url        string
		method     string
		statusCode int
	}

	testCases := []testCase{{
		url:        "/login",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/login",
		method:     http.MethodPost,
		statusCode: http.StatusBadRequest,
	}, {
		url:        "/signup",
		method:     http.MethodPost,
		statusCode: http.StatusOK,
	}, {
		url:        "/signup",
		method:     http.MethodGet,
		statusCode: http.StatusBadRequest,
	}, {
		url:        "/api/health",
		method:     http.MethodGet,
		statusCode: http.StatusOK,
	}, {
		url:        "/kek",
		method:     http.MethodGet,
		statusCode: http.StatusNotFound,
	}}

	handle := MakeMux()

	for _, tc := range testCases {
		resp := utils.MakeRequest("", tc.method, handle, nil)
		assert.Equal(t, tc.statusCode, resp.StatusCode)
	}
}

// those tests are not prarallel because os.Setenv sets env var globally, across
// all goroutines, which messess up test cases when I don't expect env vars to
// be set

func TestInitShouldReturnConfigWhenEnvVarsAreSet(t *testing.T) {
	type environmentVariable struct {
		key   string
		value string
	}

	ctx := context.TODO()
	username := "test_user"
	password := "kekekeke"
	database := "testdb"

	str := wait.ForLog("database system is ready to accept connections").
		WithOccurrence(2).
		WithStartupTimeout(5 * time.Second)
	opts := []tc.ContainerCustomizer{
		tc.WithWaitStrategy(str),
		postgres.WithUsername(username),
		postgres.WithPassword(password),
		postgres.WithDatabase(database),
	}

	cont, err := postgres.Run(ctx, db.Image, opts...)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, cont.Terminate(ctx))
	})

	host, err := cont.ContainerIP(ctx)
	require.NoError(t, err)

	envVars := []environmentVariable{
		{"PGPASSWORD", password},
		{"PGUSER", username},
		{"PGDATABASE", database},
		{"PGHOST", host},
	}

	for _, envVar := range envVars {
		require.NoError(t, os.Setenv(envVar.key, envVar.value))
	}

	t.Cleanup(func() {
		for _, envVar := range envVars {
			require.NoError(t, os.Unsetenv(envVar.key))
		}
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

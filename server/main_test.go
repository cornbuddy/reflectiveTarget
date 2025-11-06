package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// those tests are not prarallel because os.Setenv sets env var globally, across
// all goroutines, which messess up test cases when I don't expect env vars to
// be set

func TestInitShouldReturnConfigWhenEnvVarsAreSet(t *testing.T) {
	type environmentVariable struct {
		key   string
		value string
	}

	envVars := []environmentVariable{
		{"PGPASSWORD", "kek"},
		{"PGUSER", "kek"},
		{"PGDATABASE", "kek"},
		{"PGHOST", "kek"},
	}

	for _, envVar := range envVars {
		os.Setenv(envVar.key, envVar.value)
	}

	t.Cleanup(func() {
		for _, envVar := range envVars {
			os.Unsetenv(envVar.key)
		}
	})

	config, err := Init()
	require.NoError(t, err)
	require.NotEmpty(t, config)
	assert.NoError(t, config.DB.Ping())
}

func TestInitShouldReturnErrorWhenRequiredEnvVarsAreNotSet(t *testing.T) {
	_, err := Init()
	require.ErrorIs(t, err, ErrNoEnvVar)
}

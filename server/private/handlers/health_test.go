package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
)

func TestHealthHandlerShouldSucceedWhenDbWorks(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	cleanup, db, err := testdb.SetupTestDb(ctx, t)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, cleanup())
	})

	route := HealthRouter{DB: db}
	res := makeRequest(http.MethodGet, "/healthz", route.Get)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	t.Cleanup(func() { res.Body.Close() })
	assert.NoError(t, err)

	var got HealthResponse
	err = json.Unmarshal(data, &got)
	assert.NoError(t, err)
	assert.True(t, got.Connected)
	assert.GreaterOrEqual(t, got.Connections, 1)
}

func TestHealthHandlerShouldFailWhenDbDoesntWork(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	cleanup, db, err := testdb.SetupTestDb(ctx, t)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, cleanup())
	})

	db.Close()
	route := HealthRouter{DB: db}
	res := makeRequest(http.MethodGet, "/healthz", route.Get)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	t.Cleanup(func() { res.Body.Close() })
	assert.NoError(t, err)

	var got HealthResponse
	err = json.Unmarshal(data, &got)
	assert.NoError(t, err)
	assert.False(t, got.Connected)
	assert.Equal(t, got.Connections, 0)
}

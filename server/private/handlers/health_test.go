package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestHealthHandlerShouldSucceedWhenDbWorks(t *testing.T) {
	t.Parallel()

	route := healthRouter.Get
	res := utils.MakeRequest("", http.MethodGet, route, nil)
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

	cleanup, db, err := testdb.SetupTestDb(ctx, t)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, cleanup())
	})

	db.Close()
	router := HealthRouter{DB: db}
	res := utils.MakeRequest("", http.MethodGet, router.Get, nil)
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

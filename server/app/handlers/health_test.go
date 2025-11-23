package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

const healthUrl = "/health"

func TestHealthHandlerShouldSucceedWhenDbWorks(t *testing.T) {
	t.Parallel()

	res := utils.MakeRequest("", http.MethodGet, healthUrl, api, nil)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	t.Cleanup(func() { res.Body.Close() })
	assert.NoError(t, err)

	var got HealthResponse
	err = json.Unmarshal(data, &got)
	assert.NoError(t, err)
	assert.True(t, got.DbConnected)
	assert.GreaterOrEqual(t, got.DbConnections, 1)
}

func TestHealthHandlerShouldFailWhenDbDoesntWork(t *testing.T) {
	t.Parallel()

	cleanup, db, err := utils.SetupTestDb(ctx)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, cleanup())
	})

	db.Close()
	api := ApiRouter{DB: db}.Routes().ServeHTTP
	res := utils.MakeRequest("", http.MethodGet, healthUrl, api, nil)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	t.Cleanup(func() { res.Body.Close() })
	assert.NoError(t, err)

	var got HealthResponse
	err = json.Unmarshal(data, &got)
	assert.NoError(t, err)
	assert.False(t, got.DbConnected)
	assert.Equal(t, got.DbConnections, 0)
}

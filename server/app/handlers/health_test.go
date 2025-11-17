package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
	testdb "github.com/cornbuddy/reflectiveTarget/server/test/db"
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

	var got model.HealthResponse
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
	api := ApiRouter{DB: db}.Routes().ServeHTTP
	res := utils.MakeRequest("", http.MethodGet, healthUrl, api, nil)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	t.Cleanup(func() { res.Body.Close() })
	assert.NoError(t, err)

	var got model.HealthResponse
	err = json.Unmarshal(data, &got)
	assert.NoError(t, err)
	assert.False(t, got.Connected)
	assert.Equal(t, got.Connections, 0)
}

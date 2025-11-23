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
	assert.True(t, got.CacheConnected)
	assert.GreaterOrEqual(t, got.DbConnections, 1)
	assert.GreaterOrEqual(t, got.CacheConnections, 1)
}

func TestHealthHandlerShouldFailWhenDbsDontWork(t *testing.T) {
	t.Parallel()

	dbCleanup, db, err := utils.SetupTestDb(ctx)
	assert.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, dbCleanup())
	})

	cacheCleanup, cache, err := utils.SetupCache(ctx)
	assert.NoError(t, err)

	t.Cleanup(func() {
		cacheCleanup()
	})

	db.Close()
	cache.Close()
	api := ApiRouter{
		DB:    db,
		Cache: cache,
	}.Routes().ServeHTTP
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
	assert.False(t, got.CacheConnected)
	assert.Equal(t, got.DbConnections, 0)
	assert.Equal(t, got.CacheConnections, 0)
}

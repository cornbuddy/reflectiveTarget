package handlers

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

const healthUrl = "/api/health"

func TestHealthHandlerShouldSucceedWhenDbWorks(t *testing.T) {
	t.Parallel()

	res := utils.MakeRequest("", http.MethodGet, healthUrl, router, nil)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
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

	config := makeTestConfig(ctx, db, cache)
	router := NewRouter(config).ServeHTTP
	res := utils.MakeRequest("", http.MethodGet, healthUrl, router, nil)
	assert.NotNil(t, res)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

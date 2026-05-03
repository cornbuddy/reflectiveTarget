package handler

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

const healthUrl = "/api/health"

func TestHealthHandlerShouldSucceedWhenDbWorks(t *testing.T) {
	t.Parallel()

	res, _, err := utils.MakeRequest("", http.MethodGet, healthUrl, router, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHealthHandlerShouldFailWhenDbsDontWork(t *testing.T) {
	t.Parallel()

	dbCleanup, db, err := utils.StartDB(ctx)
	assert.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, dbCleanup())
	})

	cacheCleanup, cache, err := utils.SetupCache(ctx)
	assert.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, cacheCleanup())
	})

	db.Close()
	cache.Close()

	config := makeTestConfig(db, cache)
	router := MakeHandler(config).ServeHTTP
	res, _, err := utils.MakeRequest("", http.MethodGet, healthUrl, router, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

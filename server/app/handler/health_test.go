package handler_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

const healthUrl = "/api/health"

func TestHealthRouteShouldHaveProperContentType(t *testing.T) {
	t.Parallel()

	resp, _, err := utils.MakeRequest(
		ctAppJson, http.MethodGet, healthUrl, router, nil,
	)
	require.NoError(t, err)
	assert.JSONEq(t, ctAppJson, resp.Header.Get("Content-Type"))
}

func TestHealthHandlerShouldSucceedWhenDbWorks(t *testing.T) {
	t.Parallel()

	res, _, err := utils.MakeRequest("", http.MethodGet, healthUrl, router, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHealthHandlerShouldFailWhenDbsDontWork(t *testing.T) {
	t.Parallel()

	dbCleanup, db, err := utils.StartDB(ctx)
	require.NoError(t, err)

	cacheCleanup, cache, err := utils.SetupCache(ctx)
	require.NoError(t, err)

	// connections will be closed as side effect
	require.NoError(t, dbCleanup())
	require.NoError(t, cacheCleanup())

	// let's build handler with broken dependencies
	config := makeTestConfig(db, cache)
	router := handler.MakeHandler(config).ServeHTTP
	res, _, err := utils.MakeRequest("", http.MethodGet, healthUrl, router, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}

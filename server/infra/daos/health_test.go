package daos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestHealthShouldReturnErrorWhenDependenciesDontWork(t *testing.T) {
	t.Parallel()

	cleanUpDb, db, err := testutils.StartDB(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, cleanUpDb())
	})

	cleanUpCache, cache, err := testutils.SetupCache(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		cleanUpCache()
	})

	require.NoError(t, db.Close())
	require.NoError(t, cache.Close())

	health := HealthDao{
		DB:    db,
		Cache: cache,
	}
	status := health.CheckHealth(ctx)
	assert.False(t, status.CacheConnected)
	assert.False(t, status.DbConnected)
	assert.Equal(t, 0, status.CacheConnections)
	assert.Equal(t, 0, status.DbConnections)
}

func TestHealthShouldReturnNilWhenDependenciesWorksFine(t *testing.T) {
	t.Parallel()

	status := health.CheckHealth(ctx)
	assert.True(t, status.CacheConnected)
	assert.True(t, status.DbConnected)
	assert.GreaterOrEqual(t, status.CacheConnections, 1)
	assert.GreaterOrEqual(t, status.DbConnections, 1)
}

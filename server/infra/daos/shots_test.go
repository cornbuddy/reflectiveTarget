package daos_test

import (
	"database/sql"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestShotsDaoListShouldReturnEmptyListWhenNoShotsForTarget(t *testing.T) {
	t.Parallel()

	var targetID valueobjects.ID
	q := "INSERT INTO targets (name, owner_id) VALUES ($1, $2) RETURNING id"
	name := utils.MakeRandomString(5)
	require.NoError(t, db.QueryRowContext(ctx, q, name, user.ID).Scan(&targetID))

	got, err := shotsDao.List(ctx, targetID)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestShotsDaoListShouldReturnNotFoundErrorWhenNoSuchTarget(t *testing.T) {
	t.Parallel()

	got, err := shotsDao.List(ctx, 69)
	require.ErrorIs(t, err, sql.ErrNoRows)
	assert.Nil(t, got)
}

func TestShotsDaoShouldListShotsForTarget(t *testing.T) {
	t.Parallel()

	got, err := shotsDao.List(ctx, target.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestShotsDaoShouldSaveShots(t *testing.T) {
	t.Parallel()

	shooter := "kekekeke"
	shot := valueobjects.Shot{
		X: rand.IntN(101),
		Y: rand.IntN(101),
	}
	shots := valueobjects.Shots{shot}
	require.NoError(t, shotsDao.Save(ctx, shooter, target.ID, shots))

	rows, err := db.QueryContext(ctx,
		"SELECT id FROM shots WHERE x = $1 AND y = $2",
		shot.X, shot.Y,
	)

	require.NoError(t, err)
	require.NoError(t, rows.Err())

	//nolint:paralleltest
	defer rows.Close()

	assert.True(t, rows.Next(), "should save shot")
}

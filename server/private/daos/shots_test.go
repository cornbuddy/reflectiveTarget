package daos

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

func TestShotsDaoListShouldReturnEmptyListWhenNoShotsForTarget(t *testing.T) {
	t.Parallel()

	var targetID int
	q := "INSERT INTO targets (name, owner_id) VALUES ($1, $2) RETURNING id"
	require.NoError(t, db.QueryRow(q, "kek?", user.ID).Scan(&targetID))

	got, err := shotsDao.List(targetID)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestShotsDaoListShouldReturnNotFoundErrorWhenNoSuchTarget(t *testing.T) {
	t.Parallel()

	got, err := shotsDao.List(69)
	require.ErrorIs(t, err, model.ErrNotFound)
	assert.Nil(t, got)
}

func TestShotsDaoShouldListShotsForTarget(t *testing.T) {
	t.Parallel()

	got, err := shotsDao.List(target.ID)
	require.NoError(t, err)
	assert.NotEmpty(t, got)
}

func TestShotsDaoShouldSaveShots(t *testing.T) {
	t.Parallel()

	shooter := "kekekeke"
	shot := model.Shot{
		X: rand.IntN(101),
		Y: rand.IntN(101),
	}
	shots := model.Shots{shot}
	require.NoError(t, shotsDao.Save(shooter, target.ID, shots))

	res, err := db.Query(
		"SELECT * FROM shots WHERE x = $1 AND y = $2",
		shot.X, shot.Y,
	)
	require.NoError(t, err)
	assert.True(t, res.Next(), "should save shot")
}

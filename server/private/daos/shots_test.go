package daos

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

func TestShotsDaoShouldListShotsForTarget(t *testing.T) {
	t.Parallel()
	t.Fatal("not implemented")
}

func TestShotsDaoShouldSaveShots(t *testing.T) {
	t.Parallel()

	username := "shots"
	user := model.User{
		Username: username,
		Password: password,
	}
	err := userDao.Save(&user)
	require.NoError(t, err)

	var targetID int
	q := "INSERT INTO targets (name, owner_id) VALUES ($1, $2) RETURNING id"
	require.NoError(t, db.QueryRow(q, "kek?", user.ID).Scan(&targetID))

	shooter := "kekekeke"
	shot := model.Shot{
		X: rand.IntN(101),
		Y: rand.IntN(101),
	}
	shots := model.Shots{shot}
	require.NoError(t, shotsDao.Save(shooter, targetID, shots))

	res, err := db.Query(
		"SELECT * FROM shots WHERE x = $1 AND y = $2",
		shot.X, shot.Y,
	)
	require.NoError(t, err)
	assert.True(t, res.Next(), "should save shot")
}

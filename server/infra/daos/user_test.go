package daos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
)

func TestSaveShouldReturnErrorIfUsernameIsTaken(t *testing.T) {
	t.Parallel()

	user := entities.User{Username: "exists", Password: password}
	require.NoError(t, userDao.Save(ctx, &user))

	err := userDao.Save(ctx, &user)
	require.ErrorContains(t, err, "violates unique constraint")
}

func TestShouldSaveUser(t *testing.T) {
	t.Parallel()

	want := entities.User{Username: "kek2", Password: password}
	assert.NoError(t, userDao.Save(ctx, &want))

	var id int
	query := "SELECT id FROM users WHERE username = $1"
	assert.NoError(t, db.QueryRow(query, want.Username).Scan(&id))
	assert.GreaterOrEqual(t, id, 1)
	assert.Equal(t, id, want.ID)
}

func TestShouldFindUserIfExists(t *testing.T) {
	t.Parallel()

	query := "INSERT INTO users (username, hashed_password) VALUES ($1, $2)"
	username := "kek"
	_, err := db.Exec(query, username, "kek1")
	assert.NoError(t, err)

	user, err := userDao.Find(ctx, username)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)
	assert.NotEmpty(t, user.Password.Hash)
	assert.GreaterOrEqual(t, user.ID, 1)
}

func TestShouldReturnNilIfUserDoesNotExist(t *testing.T) {
	t.Parallel()

	user, err := userDao.Find(ctx, "not-exists")
	assert.NoError(t, err)
	assert.Nil(t, user)
}

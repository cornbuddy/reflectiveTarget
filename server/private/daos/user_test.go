package daos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldFindUserIfExists(t *testing.T) {
	t.Parallel()

	query := "INSERT INTO users (username, hashed_password, salt) " +
		"VALUES ($1, $2, $3)"
	username := "kek"
	_, err := db.Exec(query, username, "kek", "kek")
	assert.NoError(t, err)

	user, err := userDao.Find(username)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, 1, user.ID)
}

func TestShouldReturnNilIfUserDoesNotExist(t *testing.T) {
	t.Parallel()

	user, err := userDao.Find("not-exists")
	assert.NoError(t, err)
	assert.Nil(t, user)
}

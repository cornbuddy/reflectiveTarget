package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserShouldHashPassword(t *testing.T) {
	t.Parallel()

	username := "username"
	password := "password"
	user, err := NewUser(username, password)
	assert.NoError(t, err)
	assert.Equal(t, username, user.Username)
	assert.NotEqual(t, password, user.Password.Hash)
	assert.NotEmpty(t, user.Password)
}

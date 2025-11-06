package model

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserShouldReturnErrorWhenFormMissingRequiredKeys(t *testing.T) {
	t.Parallel()

	fields := []string{"username", "password"}
	for _, field := range fields {
		form := url.Values{}
		form.Set(field, "kek")
		user, err := NewUser(form)
		assert.ErrorIs(t, err, ErrUserFormMissingKeys)
		assert.Nil(t, user)
	}
}

func TestNewUserShouldHashPassword(t *testing.T) {
	t.Parallel()

	form := url.Values{}
	form.Set("username", "username")
	form.Set("password", "password")
	user, err := NewUser(form)
	assert.NoError(t, err)
	assert.Equal(t, form.Get("username"), user.Username)
	assert.NotEqual(t, form.Get("password"), user.Password.Hash)
	assert.NotEmpty(t, user.Password)
}

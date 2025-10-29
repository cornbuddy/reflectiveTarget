package model

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserShouldHashPassword(t *testing.T) {
	form := url.Values{}
	form.Set("username", "username")
	form.Set("password", "password")
	user, err := NewUser(form)
	assert.NoError(t, err)
	assert.Equal(t, form.Get("username"), user.Username)
}

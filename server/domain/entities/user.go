package entities

import (
	"errors"
	"net/url"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

var ErrUserFormMissingKeys = errors.New("user form missing keys")

type User struct {
	valueobjects.Password
	ID       int
	Username string
}

func NewUser(form url.Values) (*User, error) {
	missingKeys := !form.Has("username") || !form.Has("password")
	if missingKeys {
		return nil, ErrUserFormMissingKeys
	}

	password, err := valueobjects.NewPassword(form.Get("password"))
	if err != nil {
		return nil, err
	}

	return &User{
		Username: form.Get("username"),
		Password: *password,
	}, nil
}

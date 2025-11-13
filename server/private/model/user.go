package model

import (
	"net/url"
)

type User struct {
	Password
	ID       int
	Username string
}

func NewUser(form url.Values) (*User, error) {
	missingKeys := !form.Has("username") || !form.Has("password")
	if missingKeys {
		return nil, ErrUserFormMissingKeys
	}

	password, err := NewPassword(form.Get("password"))
	if err != nil {
		return nil, err
	}

	return &User{
		Username: form.Get("username"),
		Password: *password,
	}, nil
}

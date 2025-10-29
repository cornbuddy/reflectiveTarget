package model

import (
	"fmt"
	"net/url"
)

var (
	ErrUserFormMissingKeys = fmt.Errorf("user form missing keys")
)

type User struct {
	Username string
}

func NewUser(form url.Values) (*User, error) {
	missingKeys := !form.Has("username") || !form.Has("password")
	if missingKeys {
		return nil, ErrUserFormMissingKeys
	}

	return &User{Username: "username"}, nil
}

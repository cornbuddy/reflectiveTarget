package model

import (
	"net/url"
)

type User struct {
	Username string
}

func NewUser(form url.Values) (*User, error) {
	return &User{Username: "username"}, nil
}

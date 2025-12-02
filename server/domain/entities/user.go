package entities

import (
	"net/url"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type User struct {
	valueobjects.Password
	ID       int
	Username string
}

func NewUser(form url.Values) (*User, error) {
	password, err := valueobjects.NewPassword(form.Get("password"))
	if err != nil {
		return nil, err
	}

	return &User{
		Username: form.Get("username"),
		Password: *password,
	}, nil
}

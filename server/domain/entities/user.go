package entities

import (
	"fmt"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

type User struct {
	Username string
	valueobjects.Password
	valueobjects.ID
}

type Users []User

func (u *User) String() string {
	return fmt.Sprintf("ID: %d, Username: '%s'", u.ID, u.Username)
}

func NewUser(username, plaintext string) (*User, error) {
	password, err := valueobjects.NewPassword(plaintext)
	if err != nil {
		return nil, err
	}

	return &User{
		Username: username,
		Password: *password,
	}, nil
}

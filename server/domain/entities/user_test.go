package entities_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
)

func TestUserStringer(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc string
		user entities.User
		want string
	}

	testCases := []testCase{{
		"defaults",
		entities.User{},
		"{ID: 0, Username: ''}",
	}, {
		"non-defaults",
		entities.User{ID: 69, Username: "kek"},
		"{ID: 69, Username: 'kek'}",
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.user.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestNewUserShouldHashPassword(t *testing.T) {
	t.Parallel()

	username := "username"
	password := "password"
	user, err := entities.NewUser(username, password)
	require.NoError(t, err)
	assert.Equal(t, username, user.Username)
	assert.NotEqual(t, password, user.Password.Hash)
	assert.NotEmpty(t, user.Password)
}

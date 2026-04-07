package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestLoginFormValidation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		msg    string
		form   contracts.LoginForm
		want   contracts.LoginForm
		result bool
	}

	absentUsername := "do not exist"
	username := "completely random username"
	password := "default-password123@"
	user, err := entities.NewUser(username, password)
	require.NoError(t, err)
	require.NoError(t, utils.InsertUser(db, user))

	testCases := []testCase{{
		"empty fields",
		contracts.LoginForm{},
		contracts.LoginForm{
			Username: contracts.Field{Errors: contracts.Errors{ErrEmpty}},
			Password: contracts.Field{Errors: contracts.Errors{ErrEmpty}},
		},
		false,
	}, {
		"absent user without password",
		contracts.LoginForm{
			Username: contracts.Field{Value: absentUsername},
			Password: contracts.Field{},
		},
		contracts.LoginForm{
			Username: contracts.Field{
				Value:  absentUsername,
				Errors: contracts.Errors{ErrUserDoesNotExists},
			},
			Password: contracts.Field{
				Errors: contracts.Errors{ErrEmpty},
			},
		},
		false,
	}, {
		"absent user with password",
		contracts.LoginForm{
			Username: contracts.Field{Value: absentUsername},
			Password: contracts.Field{Value: "kek"},
		},
		contracts.LoginForm{
			Username: contracts.Field{
				Value:  absentUsername,
				Errors: contracts.Errors{ErrUserDoesNotExists},
			},
			Password: contracts.Field{Value: "kek"},
		},
		false,
	}, {
		"user with bad password",
		contracts.LoginForm{
			Username: contracts.Field{Value: user.Username},
			Password: contracts.Field{Value: "kek"},
		},
		contracts.LoginForm{
			Username: contracts.Field{
				Value: user.Username,
			},
			Password: contracts.Field{
				Value:  "kek",
				Errors: contracts.Errors{ErrWrongPassword},
			},
		},
		false,
	}, {
		"all good",
		contracts.LoginForm{
			Username: contracts.Field{Value: user.Username},
			Password: contracts.Field{Value: password},
		},
		contracts.LoginForm{
			Username: contracts.Field{Value: user.Username},
			Password: contracts.Field{Value: password},
		},
		true,
	}}

	validator := LoginFormValidator{userDao}
	for _, tc := range testCases {
		got := validator.Validate(ctx, &tc.form)
		assert.Equal(t, tc.want, tc.form, tc.msg)
		assert.Equal(t, tc.result, got, tc.msg)
	}
}

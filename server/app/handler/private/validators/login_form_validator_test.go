package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
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

	const absentUsername = "do not exist"
	user, err := entities.NewUser("completely random username", validPassword)
	require.NoError(t, err)
	require.NoError(t, utils.InsertUser(ctx, db, user))

	testCases := []testCase{{
		"empty fields",
		contracts.LoginForm{},
		contracts.LoginForm{
			Username: contracts.Field{Errors: contracts.Errors{validators.ErrEmpty}},
			Password: contracts.Field{Errors: contracts.Errors{validators.ErrEmpty}},
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
				Errors: contracts.Errors{validators.ErrUserDoesNotExists},
			},
			Password: contracts.Field{
				Errors: contracts.Errors{validators.ErrEmpty},
			},
		},
		false,
	}, {
		"absent user with password",
		contracts.LoginForm{
			Username: contracts.Field{Value: absentUsername},
			Password: contracts.Field{Value: weakPassword},
		},
		contracts.LoginForm{
			Username: contracts.Field{
				Value:  absentUsername,
				Errors: contracts.Errors{validators.ErrUserDoesNotExists},
			},
			Password: contracts.Field{Value: weakPassword},
		},
		false,
	}, {
		"user with bad password",
		contracts.LoginForm{
			Username: contracts.Field{Value: user.Username},
			Password: contracts.Field{Value: weakPassword},
		},
		contracts.LoginForm{
			Username: contracts.Field{
				Value: user.Username,
			},
			Password: contracts.Field{
				Value:  weakPassword,
				Errors: contracts.Errors{validators.ErrWrongPassword},
			},
		},
		false,
	}, {
		"all good",
		contracts.LoginForm{
			Username: contracts.Field{Value: user.Username},
			Password: contracts.Field{Value: validPassword},
		},
		contracts.LoginForm{
			Username: contracts.Field{Value: user.Username},
			Password: contracts.Field{Value: validPassword},
		},
		true,
	}}

	validator := validators.LoginFormValidator{userDao}
	for _, tc := range testCases {
		got := validator.Validate(ctx, &tc.form)
		assert.Equal(t, tc.want, tc.form, tc.msg)
		assert.Equal(t, tc.result, got, tc.msg)
	}
}

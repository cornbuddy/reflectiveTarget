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

func TestSignupFormValidator(t *testing.T) {
	t.Parallel()

	user, err := entities.NewUser("yet another username", validPassword)
	require.NoError(t, err)
	require.NoError(t, utils.InsertUser(db, user))

	type testCase struct {
		form  contracts.SignupForm
		want  contracts.SignupForm
		valid bool
	}

	testCases := []testCase{{
		contracts.SignupForm{
			Username:     contracts.Field{Value: ""},
			Password:     contracts.Field{Value: ""},
			Confirmation: contracts.Field{Value: ""},
		},
		contracts.SignupForm{
			Username: contracts.Field{
				Value:  "",
				Errors: contracts.Errors{validators.ErrEmpty},
			},
			Password: contracts.Field{
				Value: "",
				Errors: contracts.Errors{
					validators.ErrPasswordTooShort,
					validators.ErrPasswordDoesntContainDigits,
					validators.ErrPasswordDoesntContainSpecialChars,
				},
			},
			Confirmation: contracts.Field{
				Value: "",
			},
		},
		false,
	}, {
		contracts.SignupForm{
			Username:     contracts.Field{Value: user.Username},
			Password:     contracts.Field{Value: user.Hash},
			Confirmation: contracts.Field{Value: user.Hash},
		},
		contracts.SignupForm{
			Username: contracts.Field{
				Value:  user.Username,
				Errors: contracts.Errors{validators.ErrUserAlreadyExists},
			},
			Password: contracts.Field{
				Value: user.Hash,
			},
			Confirmation: contracts.Field{
				Value: user.Hash,
			},
		},
		false,
	}, {
		contracts.SignupForm{
			Username:     contracts.Field{Value: freeUsesrname},
			Password:     contracts.Field{Value: weakPassword},
			Confirmation: contracts.Field{Value: "not kek"},
		},
		contracts.SignupForm{
			Username: contracts.Field{
				Value: freeUsesrname,
			},
			Password: contracts.Field{
				Value: weakPassword,
				Errors: contracts.Errors{
					validators.ErrPasswordDoesntContainDigits,
					validators.ErrPasswordDoesntContainSpecialChars,
				},
			},
			Confirmation: contracts.Field{
				Value:  "not kek",
				Errors: contracts.Errors{validators.ErrPasswordsShouldMatch},
			},
		},
		false,
	}, {
		contracts.SignupForm{
			Username:     contracts.Field{Value: freeUsesrname},
			Password:     contracts.Field{Value: validPassword},
			Confirmation: contracts.Field{Value: validPassword},
		},
		contracts.SignupForm{
			Username: contracts.Field{
				Value: freeUsesrname,
			},
			Password: contracts.Field{
				Value: validPassword,
			},
			Confirmation: contracts.Field{
				Value: validPassword,
			},
		},
		true,
	}, {
		contracts.SignupForm{
			Username:     contracts.Field{Value: freeUsesrname},
			Password:     contracts.Field{Value: validPassword},
			Confirmation: contracts.Field{Value: validPassword},
		},
		contracts.SignupForm{
			Username: contracts.Field{
				Value: freeUsesrname,
			},
			Password: contracts.Field{
				Value: validPassword,
			},
			Confirmation: contracts.Field{
				Value: validPassword,
			},
		},
		true,
	}}

	validator := validators.SignupFormValidator{UserDao: userDao}
	for _, tc := range testCases {
		valid, err := validator.Validate(ctx, &tc.form)
		require.NoError(t, err)
		assert.Equal(t, tc.want, tc.form)
		assert.Equal(t, tc.valid, valid)
	}
}

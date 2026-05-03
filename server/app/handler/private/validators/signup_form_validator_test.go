package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestSignupFormValidator(t *testing.T) {
	t.Parallel()

	username := "yet another username"
	password := "default-password123@"
	user, err := entities.NewUser(username, password)
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
				Errors: contracts.Errors{ErrEmpty},
			},
			Password: contracts.Field{
				Value: "",
				Errors: contracts.Errors{
					ErrPasswordTooShort,
					ErrPasswordDoesntContainDigits,
					ErrPasswordDoesntContainSpecialChars,
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
			Password:     contracts.Field{Value: user.Password.Hash},
			Confirmation: contracts.Field{Value: user.Password.Hash},
		},
		contracts.SignupForm{
			Username: contracts.Field{
				Value:  user.Username,
				Errors: contracts.Errors{ErrUserAlreadyExists},
			},
			Password: contracts.Field{
				Value: user.Password.Hash,
			},
			Confirmation: contracts.Field{
				Value: user.Password.Hash,
			},
		},
		false,
	}, {
		contracts.SignupForm{
			Username:     contracts.Field{Value: "keker"},
			Password:     contracts.Field{Value: "kekekeke"},
			Confirmation: contracts.Field{Value: "not kek"},
		},
		contracts.SignupForm{
			Username: contracts.Field{
				Value: "keker",
			},
			Password: contracts.Field{
				Value: "kekekeke",
				Errors: contracts.Errors{
					ErrPasswordDoesntContainDigits,
					ErrPasswordDoesntContainSpecialChars,
				},
			},
			Confirmation: contracts.Field{
				Value:  "not kek",
				Errors: contracts.Errors{ErrPasswordsShouldMatch},
			},
		},
		false,
	}, {
		contracts.SignupForm{
			Username:     contracts.Field{Value: "keker"},
			Password:     contracts.Field{Value: "kekeke1@"},
			Confirmation: contracts.Field{Value: "kekeke1@"},
		},
		contracts.SignupForm{
			Username: contracts.Field{
				Value: "keker",
			},
			Password: contracts.Field{
				Value: "kekeke1@",
			},
			Confirmation: contracts.Field{
				Value: "kekeke1@",
			},
		},
		true,
	}, {
		contracts.SignupForm{
			Username:     contracts.Field{Value: "keker"},
			Password:     contracts.Field{Value: "kekeke@1"},
			Confirmation: contracts.Field{Value: "kekeke@1"},
		},
		contracts.SignupForm{
			Username: contracts.Field{
				Value: "keker",
			},
			Password: contracts.Field{
				Value: "kekeke@1",
			},
			Confirmation: contracts.Field{
				Value: "kekeke@1",
			},
		},
		true,
	}}

	validator := SignupFormValidator{UserDao: userDao}
	for _, tc := range testCases {
		valid, err := validator.Validate(ctx, &tc.form)
		require.NoError(t, err)
		assert.Equal(t, tc.want, tc.form)
		assert.Equal(t, tc.valid, valid)
	}
}

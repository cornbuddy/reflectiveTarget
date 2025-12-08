package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
)

func TestSignupFormValidator(t *testing.T) {
	t.Parallel()

	user, err := makeTestUser(userDao)
	require.NoError(t, err)

	type testCase struct {
		form forms.SignupForm
		want forms.SignupForm
	}

	testCases := []testCase{{
		forms.SignupForm{
			Username:     forms.Field{Value: ""},
			Password:     forms.Field{Value: ""},
			Confirmation: forms.Field{Value: ""},
		},
		forms.SignupForm{
			Username: forms.Field{
				Value:  "",
				Errors: forms.Errors{ErrEmpty},
			},
			Password: forms.Field{
				Value: "",
				Errors: forms.Errors{
					ErrPasswordTooShort,
					ErrPasswordDoesntContainDigits,
					ErrPasswordDoesntContainSpecialChars,
				},
			},
			Confirmation: forms.Field{
				Value: "",
			},
		},
	}, {
		forms.SignupForm{
			Username:     forms.Field{Value: user.Username},
			Password:     forms.Field{Value: user.Password.Hash},
			Confirmation: forms.Field{Value: user.Password.Hash},
		},
		forms.SignupForm{
			Username: forms.Field{
				Value:  user.Username,
				Errors: forms.Errors{ErrUserAlreadyExists},
			},
			Password: forms.Field{
				Value: user.Password.Hash,
			},
			Confirmation: forms.Field{
				Value: user.Password.Hash,
			},
		},
	}, {
		forms.SignupForm{
			Username:     forms.Field{Value: "keker"},
			Password:     forms.Field{Value: "kekekeke"},
			Confirmation: forms.Field{Value: "not kek"},
		},
		forms.SignupForm{
			Username: forms.Field{
				Value: "keker",
			},
			Password: forms.Field{
				Value: "kekekeke",
				Errors: forms.Errors{
					ErrPasswordDoesntContainDigits,
					ErrPasswordDoesntContainSpecialChars,
				},
			},
			Confirmation: forms.Field{
				Value:  "not kek",
				Errors: forms.Errors{ErrPasswordsShouldMatch},
			},
		},
	}, {
		forms.SignupForm{
			Username:     forms.Field{Value: "keker"},
			Password:     forms.Field{Value: "kekeke1@"},
			Confirmation: forms.Field{Value: "kekeke1@"},
		},
		forms.SignupForm{
			Username: forms.Field{
				Value: "keker",
			},
			Password: forms.Field{
				Value: "kekeke1@",
			},
			Confirmation: forms.Field{
				Value: "kekeke1@",
			},
		},
	}}

	validator := SignupFormValidator{UserDao: userDao}
	for _, tc := range testCases {
		validator.Validate(&tc.form)
		assert.Equal(t, tc.want, tc.form)
	}
}

package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
)

func TestSignupFormValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		form forms.SignupForm
		want forms.SignupForm
	}

	testCases := []testCase{{
		forms.SignupForm{
			Username:        forms.Field{Value: ""},
			Password:        forms.Field{Value: ""},
			ConfirmPassword: forms.Field{Value: ""},
		},
		forms.SignupForm{
			Username: forms.Field{
				Value:  "",
				Errors: forms.Errors{ErrEmpty},
			},
			Password: forms.Field{
				Value:  "",
				Errors: forms.Errors{ErrEmpty},
			},
			ConfirmPassword: forms.Field{
				Value:  "",
				Errors: forms.Errors{ErrEmpty},
			},
		},
	}}

	validator := SignupFormValidator{UserDao: userDao}
	for _, tc := range testCases {
		validator.Validate(&tc.form)
		assert.Equal(t, tc.want, tc.form)
	}

}

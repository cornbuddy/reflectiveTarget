package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
)

func TestLoginFormValidation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		form   forms.LoginForm
		want   forms.LoginForm
		result bool
	}

	user, err := makeTestUser(userDao)
	require.NoError(t, err)

	testCases := []testCase{{
		forms.LoginForm{},
		forms.LoginForm{
			Username: forms.Field{Errors: forms.Errors{ErrEmpty}},
			Password: forms.Field{Errors: forms.Errors{ErrEmpty}},
		},
		false,
	}, {
		forms.LoginForm{
			Username: forms.Field{Value: "kek"},
			Password: forms.Field{Value: "kek"},
		},
		forms.LoginForm{
			Username: forms.Field{
				Value:  "kek",
				Errors: forms.Errors{ErrUserDoesNotExists},
			},
			Password: forms.Field{Value: "kek"},
		},
		false,
	}, {
		forms.LoginForm{
			Username: forms.Field{Value: user.Username},
			Password: forms.Field{Value: "kek"},
		},
		forms.LoginForm{
			Username: forms.Field{
				Value:  user.Username,
				Errors: forms.Errors{ErrWrongPassword},
			},
			Password: forms.Field{Value: "kek"},
		},
		false,
	}, {
		forms.LoginForm{
			Username: forms.Field{Value: user.Username},
			Password: forms.Field{Value: defaultPassword},
		},
		forms.LoginForm{
			Username: forms.Field{Value: user.Username},
			Password: forms.Field{Value: defaultPassword},
		},
		true,
	}}

	validator := LoginFormValidator{userDao}
	for _, tc := range testCases {
		got := validator.Validate(&tc.form)
		assert.Equal(t, tc.want, tc.form)
		assert.Equal(t, tc.result, got)
	}
}

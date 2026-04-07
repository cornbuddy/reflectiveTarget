package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/forms"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestLoginFormValidation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		msg    string
		form   forms.LoginForm
		want   forms.LoginForm
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
		forms.LoginForm{},
		forms.LoginForm{
			Username: forms.Field{Errors: forms.Errors{ErrEmpty}},
			Password: forms.Field{Errors: forms.Errors{ErrEmpty}},
		},
		false,
	}, {
		"absent user without password",
		forms.LoginForm{
			Username: forms.Field{Value: absentUsername},
			Password: forms.Field{},
		},
		forms.LoginForm{
			Username: forms.Field{
				Value:  absentUsername,
				Errors: forms.Errors{ErrUserDoesNotExists},
			},
			Password: forms.Field{
				Errors: forms.Errors{ErrEmpty},
			},
		},
		false,
	}, {
		"absent user with password",
		forms.LoginForm{
			Username: forms.Field{Value: absentUsername},
			Password: forms.Field{Value: "kek"},
		},
		forms.LoginForm{
			Username: forms.Field{
				Value:  absentUsername,
				Errors: forms.Errors{ErrUserDoesNotExists},
			},
			Password: forms.Field{Value: "kek"},
		},
		false,
	}, {
		"user with bad password",
		forms.LoginForm{
			Username: forms.Field{Value: user.Username},
			Password: forms.Field{Value: "kek"},
		},
		forms.LoginForm{
			Username: forms.Field{
				Value: user.Username,
			},
			Password: forms.Field{
				Value:  "kek",
				Errors: forms.Errors{ErrWrongPassword},
			},
		},
		false,
	}, {
		"all good",
		forms.LoginForm{
			Username: forms.Field{Value: user.Username},
			Password: forms.Field{Value: password},
		},
		forms.LoginForm{
			Username: forms.Field{Value: user.Username},
			Password: forms.Field{Value: password},
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

package forms_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
)

func TestLoginForm(t *testing.T) {
	t.Parallel()

	type testCase struct {
		httpForm url.Values
		wantForm forms.LoginForm
	}

	testCases := []testCase{{
		url.Values{},
		forms.LoginForm{},
	}, {
		url.Values{"username": []string{"kek"}},
		forms.LoginForm{
			Username: forms.Field{Value: "kek"},
		},
	}, {
		url.Values{"password": []string{"pass"}},
		forms.LoginForm{
			Password: forms.Field{Value: "pass"},
		},
	}, {
		url.Values{
			"username": []string{"kek"},
			"password": []string{"pass"},
		},
		forms.LoginForm{
			Username: forms.Field{Value: "kek"},
			Password: forms.Field{Value: "pass"},
		},
	}}

	for _, tc := range testCases {
		gotForm := forms.NewLoginForm(tc.httpForm)
		assert.Equal(t, tc.wantForm, gotForm)
	}
}

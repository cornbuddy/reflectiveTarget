package forms_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/forms"
)

func TestSignupForm(t *testing.T) {
	t.Parallel()

	type testCase struct {
		httpForm url.Values
		wantForm forms.SignupForm
	}

	testCases := []testCase{{
		url.Values{},
		forms.SignupForm{},
	}, {
		url.Values{"username": []string{"kek"}},
		forms.SignupForm{
			Username: forms.Field{Value: "kek"},
		},
	}, {
		url.Values{"password": []string{"pass"}},
		forms.SignupForm{
			Password: forms.Field{Value: "pass"},
		},
	}, {
		url.Values{"confirmation": []string{"pass"}},
		forms.SignupForm{
			Confirmation: forms.Field{Value: "pass"},
		},
	}, {
		url.Values{
			"username":     []string{"kek"},
			"password":     []string{"pass"},
			"confirmation": []string{"pass"},
		},
		forms.SignupForm{
			Username:     forms.Field{Value: "kek"},
			Password:     forms.Field{Value: "pass"},
			Confirmation: forms.Field{Value: "pass"},
		},
	}}

	for _, tc := range testCases {
		gotForm := forms.NewSignupForm(tc.httpForm)
		assert.Equal(t, tc.wantForm, gotForm)
	}
}

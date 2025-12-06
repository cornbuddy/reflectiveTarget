package forms_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
)

func TestShouldCreateProperInstanceOfSignupForm(t *testing.T) {
	t.Parallel()

	type testCase struct {
		httpForm url.Values
		wantForm forms.SignupForm
	}

	testCases := []testCase{{
		httpForm: url.Values{},
		wantForm: forms.SignupForm{},
	}, {
		httpForm: url.Values{"username": []string{"kek"}},
		wantForm: forms.SignupForm{
			Username: forms.Field{Value: "kek"},
		},
	}, {
		httpForm: url.Values{"password": []string{"pass"}},
		wantForm: forms.SignupForm{
			Password: forms.Field{Value: "pass"},
		},
	}, {
		httpForm: url.Values{"confirmation": []string{"pass"}},
		wantForm: forms.SignupForm{
			Confirmation: forms.Field{Value: "pass"},
		},
	}, {
		httpForm: url.Values{
			"username":     []string{"kek"},
			"password":     []string{"pass"},
			"confirmation": []string{"pass"},
		},
		wantForm: forms.SignupForm{
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

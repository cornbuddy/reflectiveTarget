package formdata_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/formdata"
)

func TestShouldCreateProperInstanceOfSignupForm(t *testing.T) {
	t.Parallel()

	type testCase struct {
		httpForm url.Values
		wantForm formdata.SignupForm
	}

	testCases := []testCase{{
		httpForm: url.Values{},
		wantForm: formdata.SignupForm{},
	}, {
		httpForm: url.Values{"username": []string{"kek"}},
		wantForm: formdata.SignupForm{Username: "kek"},
	}, {
		httpForm: url.Values{"password": []string{"pass"}},
		wantForm: formdata.SignupForm{Password: "pass"},
	}, {
		httpForm: url.Values{"repeated-password": []string{"pass"}},
		wantForm: formdata.SignupForm{RepeatedPassword: "pass"},
	}, {
		httpForm: url.Values{
			"username":          []string{"kek"},
			"password":          []string{"pass"},
			"repeated-password": []string{"pass"},
		},
		wantForm: formdata.SignupForm{
			Username:         "kek",
			Password:         "pass",
			RepeatedPassword: "pass",
		},
	}}

	for _, tc := range testCases {
		gotForm := formdata.NewSignupForm(tc.httpForm)
		assert.Equal(t, tc.wantForm, gotForm)
	}
}

package contracts_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

func TestSignupForm(t *testing.T) {
	t.Parallel()

	type testCase struct {
		httpForm url.Values
		wantForm contracts.SignupForm
	}

	testCases := []testCase{{
		url.Values{},
		contracts.SignupForm{},
	}, {
		url.Values{"username": []string{"kek"}},
		contracts.SignupForm{
			Username: contracts.Field{Value: "kek"},
		},
	}, {
		url.Values{"password": []string{"pass"}},
		contracts.SignupForm{
			Password: contracts.Field{Value: "pass"},
		},
	}, {
		url.Values{"confirmation": []string{"pass"}},
		contracts.SignupForm{
			Confirmation: contracts.Field{Value: "pass"},
		},
	}, {
		url.Values{
			"username":     []string{"kek"},
			"password":     []string{"pass"},
			"confirmation": []string{"pass"},
		},
		contracts.SignupForm{
			Username:     contracts.Field{Value: "kek"},
			Password:     contracts.Field{Value: "pass"},
			Confirmation: contracts.Field{Value: "pass"},
		},
	}}

	for _, tc := range testCases {
		gotForm := contracts.NewSignupForm(tc.httpForm)
		assert.Equal(t, tc.wantForm, gotForm)
	}
}

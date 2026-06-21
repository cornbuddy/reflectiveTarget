package contracts_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

const formConfirmation = "confirmation"

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
		url.Values{formUsername: []string{username}},
		contracts.SignupForm{
			Username: contracts.Field{Value: username},
		},
	}, {
		url.Values{formPassword: []string{password}},
		contracts.SignupForm{
			Password: contracts.Field{Value: password},
		},
	}, {
		url.Values{formConfirmation: []string{password}},
		contracts.SignupForm{
			Confirmation: contracts.Field{Value: password},
		},
	}, {
		url.Values{
			formUsername:     []string{username},
			formPassword:     []string{password},
			formConfirmation: []string{password},
		},
		contracts.SignupForm{
			Username:     contracts.Field{Value: username},
			Password:     contracts.Field{Value: password},
			Confirmation: contracts.Field{Value: password},
		},
	}}

	for _, tc := range testCases {
		gotForm := contracts.NewSignupForm(tc.httpForm)
		assert.Equal(t, tc.wantForm, gotForm)
	}
}

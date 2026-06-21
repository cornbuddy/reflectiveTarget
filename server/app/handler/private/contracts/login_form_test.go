package contracts_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

func TestLoginForm(t *testing.T) {
	t.Parallel()

	type testCase struct {
		httpForm url.Values
		wantForm contracts.LoginForm
	}

	testCases := []testCase{{
		url.Values{},
		contracts.LoginForm{},
	}, {
		url.Values{formUsername: []string{username}},
		contracts.LoginForm{
			Username: contracts.Field{Value: username},
		},
	}, {
		url.Values{formPassword: []string{password}},
		contracts.LoginForm{
			Password: contracts.Field{Value: password},
		},
	}, {
		url.Values{
			formUsername: []string{username},
			formPassword: []string{password},
		},
		contracts.LoginForm{
			Username: contracts.Field{Value: username},
			Password: contracts.Field{Value: password},
		},
	}}

	for _, tc := range testCases {
		gotForm := contracts.NewLoginForm(tc.httpForm)
		assert.Equal(t, tc.wantForm, gotForm)
	}
}

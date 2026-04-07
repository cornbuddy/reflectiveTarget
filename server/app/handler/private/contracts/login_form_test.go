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
		url.Values{"username": []string{"kek"}},
		contracts.LoginForm{
			Username: contracts.Field{Value: "kek"},
		},
	}, {
		url.Values{"password": []string{"pass"}},
		contracts.LoginForm{
			Password: contracts.Field{Value: "pass"},
		},
	}, {
		url.Values{
			"username": []string{"kek"},
			"password": []string{"pass"},
		},
		contracts.LoginForm{
			Username: contracts.Field{Value: "kek"},
			Password: contracts.Field{Value: "pass"},
		},
	}}

	for _, tc := range testCases {
		gotForm := contracts.NewLoginForm(tc.httpForm)
		assert.Equal(t, tc.wantForm, gotForm)
	}
}

package entities

import (
	"net/url"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserFieldsShouldBePopulatedAccordingly(t *testing.T) {
	t.Parallel()

	type testCase struct {
		form url.Values
		user User
	}

	testCases := []testCase{{
		form: url.Values{},
		user: User{},
	}, {
		form: url.Values{"username": []string{"kek"}},
		user: User{Username: "kek"},
	}, {
		form: url.Values{"password": []string{"kek"}},
		user: User{
			Password: valueobjects.Password{
				Hash:      "not kek",
				Plaintext: "kek",
			},
		},
	}, {
		form: url.Values{
			"password": []string{"kek"},
			"username": []string{"kek"},
		},
		user: User{
			Username: "kek",
			Password: valueobjects.Password{
				Hash:      "not kek",
				Plaintext: "kek",
			},
		},
	}}

	for _, tc := range testCases {
		user, err := NewUser(tc.form)
		require.NoError(t, err)
		assert.Equal(t, tc.user, user)
	}

}

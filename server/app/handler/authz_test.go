package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestLogoutShouldUpdateSessionCookie(t *testing.T) {
	t.Parallel()

	const url = "/logout"

	ct := "application/x-www-form-urlencoded"
	res, body, err := utils.MakeRequest(ct, http.MethodGet, url, router, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, res.StatusCode)
	assert.Equal(t, "/", res.Header.Get("Location"))
	assertAuthenticationStatusIsChanged(t, res,
		"logout should update session cookie",
	)
	assert.Contains(t, body, "Logout succeeded")
}

func TestLoginShouldFailWhenSomethingIsWrong(t *testing.T) {
	t.Parallel()

	type testCase struct {
		message    string
		statusCode int
		body       io.Reader
	}

	const url = "/login"
	wrongPasswordUser, err := makeTestUser(db)
	require.NoError(t, err)
	require.NotNil(t, wrongPasswordUser)

	testCases := []testCase{{
		message:    "user does not exist",
		statusCode: http.StatusUnauthorized,
		body: strings.NewReader(
			fmt.Sprintf("username=%s&password=%s", "absent", "kek"),
		),
	}, {
		message:    "wrong password",
		statusCode: http.StatusUnauthorized,
		body: strings.NewReader(
			fmt.Sprintf(
				"username=%s&password=%s",
				wrongPasswordUser.Username,
				"wrong-password",
			),
		),
	}}

	ct := "application/x-www-form-urlencoded"
	for _, tc := range testCases {
		res, body, err := utils.MakeRequest(ct, http.MethodPost, url, router, tc.body)
		require.NoError(t, err)
		assert.Equal(t, tc.statusCode, res.StatusCode)
		assert.Contains(t, body, tc.message)
	}
}

func TestLoginShouldSetSessionCookieOnSuccess(t *testing.T) {
	t.Parallel()

	const url = "/login"

	user, err := makeTestUser(db)
	require.NoError(t, err)

	ct := "application/x-www-form-urlencoded"
	b := strings.NewReader(
		fmt.Sprintf(
			"username=%s&password=%s",
			user.Username, defaultPassword,
		),
	)
	res, body, err := utils.MakeRequest(ct, http.MethodPost, url, router, b)
	require.NoError(t, err)
	assertAuthenticationStatusIsChanged(t, res,
		"login should set session cookie",
	)
	assert.Contains(t, body, "Login succeeded")
}

func TestShouldRegisterNewUserWhenCredentialsAreValid(t *testing.T) {
	t.Parallel()

	type testCase struct {
		message    string
		statusCode int
		body       io.Reader
	}

	const url = "/signup"

	testCases := []testCase{{
		message:    "User successfully created",
		statusCode: http.StatusSeeOther,
		body: strings.NewReader(
			fmt.Sprintf(
				"username=%s&password=%s&confirmation=%s",
				"kek", defaultPassword, defaultPassword,
			),
		),
	}, {
		message:    "user already exists",
		statusCode: http.StatusBadRequest,
		body: strings.NewReader(
			fmt.Sprintf(
				"username=%s&password=%s&confirmation=%s",
				"kek", defaultPassword, defaultPassword,
			),
		),
	}, {
		message:    "password should be at least 8 characters long",
		statusCode: http.StatusBadRequest,
		body: strings.NewReader(
			fmt.Sprintf("username=%s", "kek"),
		),
	}}

	ct := "application/x-www-form-urlencoded"
	for _, tc := range testCases {
		res, body, err := utils.MakeRequest(ct, http.MethodPost, url, router, tc.body)
		require.NoError(t, err)
		assert.Equal(t, tc.statusCode, res.StatusCode)
		assert.Contains(t, body, tc.message)

		isSucceed := res.StatusCode == http.StatusSeeOther
		if isSucceed {
			assertAuthenticationStatusIsChanged(t, res, tc.message)
		}
	}
}

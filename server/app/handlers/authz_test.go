package handlers

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
	res := utils.MakeRequest(ct, http.MethodGet, url, router, nil)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusSeeOther, res.StatusCode)
	assert.Equal(t, "/", res.Header.Get("Location"))
	assertAuthenticationStatusIsChanged(t, res)

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	t.Cleanup(func() { res.Body.Close() })

	gotBody := string(data)
	assert.Contains(t, gotBody, "Logout succeeded")
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

	for _, tc := range testCases {
		ct := "application/x-www-form-urlencoded"
		body := tc.body
		res := utils.MakeRequest(ct, http.MethodPost, url, router, body)
		require.NotNil(t, res)
		assert.Equal(t, tc.statusCode, res.StatusCode)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		t.Cleanup(func() { res.Body.Close() })

		gotBody := string(data)
		assert.Contains(t, gotBody, tc.message)
	}
}

func TestLoginShouldSetSessionCookieOnSuccess(t *testing.T) {
	t.Parallel()

	const url = "/login"

	user, err := makeTestUser(db)
	require.NoError(t, err)

	ct := "application/x-www-form-urlencoded"
	body := strings.NewReader(
		fmt.Sprintf(
			"username=%s&password=%s",
			user.Username, defaultPassword,
		),
	)
	res := utils.MakeRequest(ct, http.MethodPost, url, router, body)
	require.NotNil(t, res)
	assertAuthenticationStatusIsChanged(t, res)

	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	t.Cleanup(func() { res.Body.Close() })

	gotBody := string(data)
	assert.Contains(t, gotBody, "Login succeeded")
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

	for _, tc := range testCases {
		ct := "application/x-www-form-urlencoded"
		body := tc.body
		res := utils.MakeRequest(ct, http.MethodPost, url, router, body)
		require.NotNil(t, res)
		assert.Equal(t, tc.statusCode, res.StatusCode)

		isSucceed := res.StatusCode == http.StatusSeeOther
		if isSucceed {
			assertAuthenticationStatusIsChanged(t, res)
		}

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		t.Cleanup(func() { res.Body.Close() })

		gotBody := string(data)
		assert.Contains(t, gotBody, tc.message)
	}
}

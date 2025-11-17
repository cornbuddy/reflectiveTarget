package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/model/entities"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestLoginShouldFailWhenSomethingIsWrong(t *testing.T) {
	t.Parallel()

	type testCase struct {
		message    string
		statusCode int
		body       io.Reader
	}

	const url = "/login"
	wrongPasswordUser, err := makeTestUser(userDao)
	require.NoError(t, err)
	require.NotNil(t, wrongPasswordUser)

	testCases := []testCase{{
		message:    "User does not exist",
		statusCode: http.StatusUnauthorized,
		body: strings.NewReader(
			fmt.Sprintf("username=%s&password=%s", "absent", "kek"),
		),
	}, {
		message:    "Password is wrong",
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
		res := utils.MakeRequest(ct, http.MethodPost, url, views, body)
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

	username := "test-login"
	password := "kek"
	pwd, err := entities.NewPassword(password)
	require.NoError(t, err)

	user := entities.User{
		Username: username,
		Password: *pwd,
	}
	require.NoError(t, userDao.Save(&user))

	ct := "application/x-www-form-urlencoded"
	body := strings.NewReader(
		fmt.Sprintf("username=%s&password=%s", username, password),
	)
	res := utils.MakeRequest(ct, http.MethodPost, url, views, body)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusSeeOther, res.StatusCode)
	assert.Equal(t, "/", res.Header.Get("Location"))
	assertSessionCookieIsSet(t, res)

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
			fmt.Sprintf("username=%s&password=%s", "kek", "kek"),
		),
	}, {
		message:    "User already exists",
		statusCode: http.StatusConflict,
		body: strings.NewReader(
			fmt.Sprintf("username=%s&password=%s", "kek", "kek"),
		),
	}, {
		message:    "missing keys",
		statusCode: http.StatusBadRequest,
		body: strings.NewReader(
			fmt.Sprintf("username=%s", "kek"),
		),
	}}

	for _, tc := range testCases {
		ct := "application/x-www-form-urlencoded"
		body := tc.body
		res := utils.MakeRequest(ct, http.MethodPost, url, views, body)
		require.NotNil(t, res)
		assert.Equal(t, tc.statusCode, res.StatusCode)

		isSucceed := res.StatusCode == http.StatusSeeOther
		if isSucceed {
			url := res.Header.Get("Location")
			assert.Equal(t, "/", url)
			assertSessionCookieIsSet(t, res)
		}

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		t.Cleanup(func() { res.Body.Close() })

		gotBody := string(data)
		assert.Contains(t, gotBody, tc.message)
	}
}

func TestAuthzHandlerShouldRenderFormsOnGet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		url      string
		contains string
	}

	testCases := []testCase{{
		url:      "/signup",
		contains: "Signup",
	}, {
		url:      "/login",
		contains: "Login",
	}}

	for _, tc := range testCases {
		get := http.MethodGet
		res := utils.MakeRequest("", get, tc.url, views, nil)
		ct := "text/html; charset=utf-8"
		assert.Equal(t, ct, res.Header.Get("Content-Type"))
		assert.Equal(t, http.StatusOK, res.StatusCode)

		data, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		t.Cleanup(func() { res.Body.Close() })

		body := string(data)
		assert.Contains(t, body, tc.contains)
		assert.Contains(t, body, "</form>")
	}
}

package handlers

import (
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/model"
)

var defaultPassword = "default-password"

func makeTestUser(dao daos.UserDao) (*model.User, error) {
	pwd, err := model.NewPassword(defaultPassword)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Username: makeRandomString(10),
		Password: *pwd,
	}
	if err := dao.Save(user); err != nil {
		return nil, err
	}

	return &user, nil
}

func makeRandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	result := make([]rune, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func assertSessionCookieIsSet(t *testing.T, resp *http.Response) {
	cookies := resp.Cookies()
	require.NotEmpty(t, cookies)

	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == SessionCookieName {
			sessionCookie = cookie
			break
		}
	}

	require.NotNil(t, sessionCookie)

	// substracting 1 second because time.Now add ms to the time, while
	// cookie doesn't count that strictly
	month := time.Now().AddDate(0, 1, 0).Add(-1 * time.Second).UTC()
	assert.True(t, sessionCookie.Expires.After(month))
	assert.NoError(t, uuid.Validate(sessionCookie.Value))
}

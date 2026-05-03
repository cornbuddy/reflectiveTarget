package handler

import (
	"database/sql"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appconst "github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

const defaultPassword = "default-password123@"

// TODO: remove
func makeTestTarget(db *sql.DB, userID int) (int, error) {
	var targetID int
	q := "INSERT INTO targets (name, owner_id) VALUES($1, $2) " +
		"RETURNING id"
	if err := db.QueryRow(q, "kek?", userID).Scan(&targetID); err != nil {
		return 0, err
	}

	return targetID, nil
}

// TODO: remove
func makeTestUser(db *sql.DB) (*entities.User, error) {
	username := utils.MakeRandomString(10)
	user, err := entities.NewUser(username, defaultPassword)
	if err != nil {
		return nil, err
	}

	if err := utils.InsertUser(db, user); err != nil {
		return nil, err
	}

	return user, nil
}

func assertSessionCookieIsSet(t *testing.T, resp *http.Response, msg string) {
	cookies := resp.Cookies()
	require.NotEmpty(t, cookies, msg)

	var sessionCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == appconst.SessionCookieName {
			sessionCookie = cookie
			break
		}
	}

	require.NotNil(t, sessionCookie, msg)

	duration := time.Now().
		Add(sessiondata.SessionDuration).
		// substracting few seconds because actual tests can happen
		// after time.Now()
		Add(-10 * time.Second)
	assert.True(t, sessionCookie.Expires.After(duration), msg)
	assert.NoError(t, uuid.Validate(sessionCookie.Value), msg)
}

func assertAuthenticationStatusIsChanged(
	t *testing.T, resp *http.Response, msg string,
) {

	assertSessionCookieIsSet(t, resp, msg)

	url := resp.Header.Get("Location")
	assert.Equal(t, "/", url, "location should be set")
}

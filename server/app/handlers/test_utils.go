package handlers

import (
	"database/sql"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appconst "github.com/cornbuddy/reflectiveTarget/server/app/constants"
	domconst "github.com/cornbuddy/reflectiveTarget/server/domain/constants"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

const defaultPassword = "default-password123@"

func makeTestTarget(db *sql.DB, userID int) (int, error) {
	var targetID int
	q := "INSERT INTO targets (name, owner_id) VALUES($1, $2) " +
		"RETURNING id"
	if err := db.QueryRow(q, "kek?", userID).Scan(&targetID); err != nil {
		return 0, err
	}

	return targetID, nil
}

func makeTestUser(db *sql.DB) (*entities.User, error) {
	username := makeRandomString(10)
	user, err := entities.NewUser(username, defaultPassword)
	if err != nil {
		return nil, err
	}

	if err := utils.InsertUser(db, user); err != nil {
		return nil, err
	}

	return user, nil
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
		if cookie.Name == appconst.SessionCookieName {
			sessionCookie = cookie
			break
		}
	}

	require.NotNil(t, sessionCookie)

	// substracting few seconds because actual tests can happen after
	// time.Now()
	month := time.Now().Add(domconst.SessionDuration).Add(-10 * time.Second)
	assert.True(t, sessionCookie.Expires.After(month))
	assert.NoError(t, uuid.Validate(sessionCookie.Value))
}

func assertAuthenticationStatusIsChanged(t *testing.T, resp *http.Response) {
	assertSessionCookieIsSet(t, resp)

	url := resp.Header.Get("Location")
	assert.Equal(t, "/", url, "location should be set")
}

package utils

import (
	"net/http"
	"time"

	"github.com/google/uuid"

	appconst "github.com/cornbuddy/reflectiveTarget/server/app/constants"
	domconst "github.com/cornbuddy/reflectiveTarget/server/domain/constants"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

func SaveSession(
	store daos.SessionStore, isAuthenticated bool, w http.ResponseWriter,
) (string, error) {

	raw, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	token := raw.String()
	if err := store.SaveSession(token, isAuthenticated); err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    appconst.SessionCookieName,
		Value:   token,
		Expires: time.Now().Add(domconst.SessionDuration),
	})

	return token, nil
}

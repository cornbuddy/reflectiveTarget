package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	appconst "github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
)

func SaveSession(
	ctx context.Context, store daos.SessionStore,
	data sessiondata.SessionData, w http.ResponseWriter,
) (string, error) {

	raw, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	token := raw.String()
	if err := store.Update(ctx, token, data); err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    appconst.SessionCookieName,
		Value:   token,
		Expires: time.Now().Add(sessiondata.SessionDuration),
	})

	return token, nil
}

package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/session"
)

func SaveSession(
	ctx context.Context, store session.Store,
	data session.Data, w http.ResponseWriter,
) (session.SessionID, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	raw := uuid.String()
	token := session.SessionID(raw)
	if err := store.Update(ctx, token, data); err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    raw,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(session.Duration),
	})

	return token, nil
}

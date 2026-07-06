package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/session"
)

func SaveSession(
	ctx context.Context, store session.Store,
	data session.Data, w http.ResponseWriter,
) (*session.SessionID, error) {
	id := session.MakeID()
	if err := store.Update(ctx, id, data); err != nil {
		return nil, err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    id.String(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(session.Duration),
	})

	return &id, nil
}

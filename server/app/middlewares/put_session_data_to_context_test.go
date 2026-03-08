package middlewares

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestPutSessionDataToCtx(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc      string
		sessionId string
		want      sessiondata.SessionData
	}

	sessionId := "kekeke"
	want := sessiondata.SessionData{
		IsAuthenticated: true,
		UserID:          69,
		Username:        "kek",
	}
	require.NoError(t, store.Update(ctx, sessionId, want))

	testCases := []testCase{{
		"should put cached data to context if any",
		sessionId,
		want,
	}, {
		"should put zero object on cache miss",
		"not found",
		sessiondata.SessionData{},
	}}

	handler := mw.PutSessionDataToContext(emptyStub).ServeHTTP
	for _, tc := range testCases {
		cookies := []*http.Cookie{{
			Name:  constants.SessionCookieName,
			Value: tc.sessionId,
		}}
		resp := utils.MakeRequestWithCookies(
			"", http.MethodGet, "/", handler, nil, cookies...,
		)
		got := sessiondata.Read(resp.Request.Context())
		assert.EqualExportedValues(t, tc.want, got, tc.desc)
	}
}

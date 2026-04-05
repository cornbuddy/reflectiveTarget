package render_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
)

type (
	assertFunc             func(assert.TestingT, any, any, ...any) bool
	renderClosure          func(render.Render) func(http.ResponseWriter)
	layoutMarkersTestCases []layoutMarkersTestCase
	layoutMarkersTestCase  struct {
		desc   string
		render render.Render
		assert func(*testing.T, io.Reader, []string)
	}
)

var (
	layoutMakrers = []string{"<!DOCTYPE html>", "</html>"}
	anonCtx       = context.TODO()
	testCases     = layoutMarkersTestCases{{
		"layout render",
		render.Layout,
		assertContainsTokens,
	}, {
		"view render",
		render.View,
		assertNotContainsTokens,
	}}
	userCtx = context.WithValue(
		context.TODO(),
		sessiondata.SessionDataCtx,
		&sessiondata.SessionData{
			IsAuthenticated: true,
			UserID:          12,
			Username:        "kek",
		},
	)
)

func (tcs *layoutMarkersTestCases) run(t *testing.T, rc renderClosure) {
	for _, tc := range *tcs {
		t.Run(tc.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			rc(tc.render)(w)
			tc.assert(t, w.Body, layoutMakrers)
		})
	}
}

func assertContainsTokens(t *testing.T, body io.Reader, tokens []string) {
	assertTokens(t, body, tokens, assert.Contains)
}

func assertNotContainsTokens(t *testing.T, body io.Reader, tokens []string) {
	assertTokens(t, body, tokens, assert.NotContains)
}

func assertTokens(
	t *testing.T, body io.Reader, tokens []string, asrt assertFunc,
) {
	raw, err := io.ReadAll(body)
	require.NoError(t, err)

	strBody := string(raw)
	for _, token := range tokens {
		asrt(t, strBody, token)
	}
}

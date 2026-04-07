package render_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type (
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
		utils.AssertContainsTokens,
	}, {
		"view render",
		render.View,
		utils.AssertNotContainsTokens,
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

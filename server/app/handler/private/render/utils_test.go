package render_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/session"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type (
	renderClosure        func(render.Render) func(http.ResponseWriter)
	fieldErrorsTestCases []struct {
		selector string
		errors   contracts.Errors
	}
	layoutMarkersTestCases []struct {
		desc   string
		render render.Render
		assert func(*testing.T, *http.Response, []string)
	}
)

var (
	layoutMakrers = []string{
		"<!doctype html>",
		"</html>",
		`<meta charset="UTF-8">`,
		`<link rel="stylesheet" href="/static/index.css">`,
		`<script src="/static/index.js"></script>`,
	}
	layoutTestCases = layoutMarkersTestCases{{
		"layout render",
		render.Layout,
		utils.AssertContainsTokens,
	}, {
		"view render",
		render.View,
		utils.AssertNotContainsTokens,
	}}
	anonCtx = context.TODO()
	userCtx = context.WithValue(
		context.TODO(),
		session.SessionDataCtx,
		&session.Data{
			IsAuthenticated: true,
			UserID:          12,
			Username:        "kek",
		},
	)
)

func (tcs *fieldErrorsTestCases) assert(t *testing.T, body io.Reader) {
	doc, err := goquery.NewDocumentFromReader(body)
	require.NoError(t, err)

	for _, tc := range *tcs {
		t.Run(tc.selector, func(t *testing.T) {
			ul := doc.Find(tc.selector)
			assert.NotEmpty(t, ul.Nodes)

			wantLen := 0
			ul.Children().Each(func(_ int, li *goquery.Selection) {
				wantLen++
				//nolint:err113
				err := errors.New(li.Text())
				assert.Contains(t, tc.errors, err)
			})

			assert.Len(t, tc.errors, wantLen)
		})
	}
}

func (tcs *layoutMarkersTestCases) assert(t *testing.T, rc renderClosure) {
	for _, tc := range *tcs {
		t.Run(tc.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			rc(tc.render)(w)
			tc.assert(t, w.Result(), layoutMakrers)
		})
	}
}

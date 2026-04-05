package render_test

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
)

func TestPageShouldContainText(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		render      render.Render
		ctx         context.Context
		data        any
		contains    []string
	}

	emptyCtx := context.TODO()
	authorizedCtx := context.WithValue(
		context.TODO(),
		sessiondata.SessionDataCtx,
		&sessiondata.SessionData{
			IsAuthenticated: true,
			UserID:          12,
			Username:        "kek",
		},
	)

	testCases := []testCase{{
		"non authorized client should see correct navbar",
		render.Layout,
		emptyCtx,
		nil,
		[]string{"Login", "Signup"},
	}, {
		"authorized client should see correct navbar",
		render.Layout,
		authorizedCtx,
		nil,
		[]string{"Targets", "Logout"},
	}, {
		"signup page should contain form",
		render.View,
		emptyCtx,
		nil,
		[]string{"Signup", "<form hx-post=\"/signup\"", "</form>"},
	}, {
		"login page should contain form",
		render.View,
		emptyCtx,
		nil,
		[]string{"Login", "<form hx-post=\"/login\"", "</form>"},
	}, {
		"targets page should allow to create new target",
		render.View,
		authorizedCtx,
		nil,
		[]string{"Targets", "New target", "href=\"/targets/new\""},
	}, {
		"targets page should have links to targets",
		render.View,
		authorizedCtx,
		map[string]any{
			"Targets": aggregations.Targets{{
				ID:   1,
				Name: "kek1?",
			}, {
				ID:   69,
				Name: "kek69?",
			}},
		},
		[]string{
			"<a href=\"/targets/1\">kek1?</a>",
			"<a href=\"/targets/69\">kek69?</a>",
		},
	}}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		tc.render(tc.ctx, w, tc.data)
		raw, err := io.ReadAll(w.Body)
		require.NoError(t, err, tc.description)

		body := string(raw)
		for _, token := range tc.contains {
			assert.Contains(t, body, token, tc.description)
		}
	}
}

type renderTestCase struct {
	desc   string
	render render.Render
	assert func(*testing.T, io.Reader, []string)
}

var renderTestCases = []renderTestCase{{
	"layout render",
	render.Layout,
	assertContainsTokens,
}, {
	"view render",
	render.View,
	assertNotContainsTokens,
}}

func TestRenderIndex(t *testing.T) {
	t.Parallel()

	for _, tc := range renderTestCases {
		t.Run(tc.desc, func(t *testing.T) {
			w := httptest.NewRecorder()
			tc.render.Index(context.TODO(), w)
			tc.assert(t, w.Body, layoutMakrers)
		})
	}
}

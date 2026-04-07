package render_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type UpdateTargetTest struct{}

func (u *UpdateTargetTest) RendersProperViewFor(t *testgroup.T) {
	type testCase struct {
		desc   string
		data   render.TargetData
		tokens []string
	}

	target := aggregations.Target{Name: "kek", ID: 69}
	testCases := []testCase{{
		"empty target",
		render.TargetData{},
		[]string{
			"New target",
			"<form",
			"hx-post=\"/target/new\"",
			"hx-trigger=\"submit\"",
			"hx-target=\"main\"",
			"</form>",
			"<canvas></canvas>",
			"<button type=\"submit\">Create</button>",
		},
	}, {
		"non empty target",
		render.TargetData{target},
		[]string{
			fmt.Sprintf("<h2>%s</h2>", target.Name),
			"<form",
			fmt.Sprintf("hx-put=\"/target/%d\"", target.ID),
			"hx-trigger=\"submit\"",
			"hx-target=\"main\"",
			"</form>",
			"<canvas></canvas>",
			"<button type=\"submit\">Update</button>",
		},
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testgroup.T) {
			w := httptest.NewRecorder()
			render.View.UpdateTarget(anonCtx, w, tc.data)
			utils.AssertContainsTokens(t.T, w.Body, tc.tokens)
		})
	}
}

func (u *UpdateTargetTest) HasProperLayoutMarkers(t *testgroup.T) {
	testCases.run(t.T, func(r render.Render) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			r.UpdateTarget(anonCtx, w, render.TargetData{})
		}
	})
}

func TestUpdateTarget(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(UpdateTargetTest))
}

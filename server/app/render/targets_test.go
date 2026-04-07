package render_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type RenderTargetsTest struct{}

func (r *RenderTargetsTest) RendersTargetLinksProperly(t *testgroup.T) {
	data := render.TargetsData{aggregations.Targets{{
		ID:   1,
		Name: "kek1?",
	}, {
		ID:   69,
		Name: "kek69?",
	}}}
	tokens := []string{}
	for _, t := range data.Targets {
		a := fmt.Sprintf("<a href=\"/targets/%d\">%s</a>", t.ID, t.Name)
		tokens = append(tokens, a)
	}

	w := httptest.NewRecorder()
	render.View.Targets(anonCtx, w, data)
	utils.AssertContainsTokens(t.T, w.Body, tokens)
}

func (r *RenderTargetsTest) AllowsToCreateNewTarget(t *testgroup.T) {
	w := httptest.NewRecorder()
	render.View.Targets(anonCtx, w, render.TargetsData{})
	tokens := []string{"Targets", "New target", "href=\"/targets/new\""}
	utils.AssertContainsTokens(t.T, w.Body, tokens)
}

func (r *RenderTargetsTest) HasProperLayoutMarkers(t *testgroup.T) {
	testCases.run(t.T, func(r render.Render) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			r.Targets(anonCtx, w, render.TargetsData{})
		}
	})
}

func TestRenderTargets(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(RenderTargetsTest))
}

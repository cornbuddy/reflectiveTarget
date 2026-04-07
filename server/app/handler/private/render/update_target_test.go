package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type UpdateTargetTest struct{}

func (u *UpdateTargetTest) ContainsTargetNameIfData(t *testgroup.T) {
	w := httptest.NewRecorder()
	target := aggregations.Target{Name: "kek"}
	render.View.UpdateTarget(anonCtx, w, render.TargetData{target})
	utils.AssertContainsTokens(t.T, w.Body, []string{target.Name})
}

func (u *UpdateTargetTest) RendersProperViewForNewTarget(t *testgroup.T) {
	w := httptest.NewRecorder()
	render.View.UpdateTarget(anonCtx, w, render.TargetData{})
	utils.AssertContainsTokens(t.T, w.Body, []string{
		"New target",
		"<form hx-post=\"/target/new\" hx-trigger=\"submit\"",
		"</form>",
	})
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

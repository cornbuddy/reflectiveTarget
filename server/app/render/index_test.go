package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"
)

type RenderIndexTest struct{}

func (r *RenderIndexTest) UserSeesButtons(t *testgroup.T) {
	w := httptest.NewRecorder()
	render.Layout.Index(userCtx, w)
	tokens := []string{"Targets", "Logout"}
	assertContainsTokens(t.T, w.Body, tokens)
}

func (r *RenderIndexTest) AnonSeesButtons(t *testgroup.T) {
	w := httptest.NewRecorder()
	render.Layout.Index(anonCtx, w)
	tokens := []string{"Login", "Signup"}
	assertContainsTokens(t.T, w.Body, tokens)
}

func (r *RenderIndexTest) HasProperLayoutMarkers(t *testgroup.T) {
	testCases.run(t.T, func(r render.Render) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			r.Index(anonCtx, w)
		}
	})
}

func TestRenderIndex(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(RenderIndexTest))
}

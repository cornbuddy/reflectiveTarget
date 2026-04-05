package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"
)

type RenderLoginTest struct{}

func (r *RenderLoginTest) HasSignupForm(t *testgroup.T) {
	w := httptest.NewRecorder()
	render.View.Signup(anonCtx, w, render.SignupData{})
	tokens := []string{"Login", "<form hx-post=\"/login\"", "</form>"}
	assertContainsTokens(t.T, w.Body, tokens)
}

func (r *RenderLoginTest) HasProperLayoutMarkers(t *testgroup.T) {
	testCases.run(t.T, func(r render.Render) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			r.Login(anonCtx, w, render.LoginData{})
		}
	})
}

func TestRenderLogin(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(RenderLoginTest))
}

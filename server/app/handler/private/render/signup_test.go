package render_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type RenderSignupTest struct {
	erroredData render.SignupData
}

func (r *RenderSignupTest) RendersValidationErrors(t *testgroup.T) {
	testCases := fieldErrorsTestCases{{
		"section#username ul.errors",
		errs,
	}, {
		"section#password ul.errors",
		errs,
	}, {
		"section#confirmation ul.errors",
		errs,
	}}

	w := httptest.NewRecorder()
	render.View.Signup(userCtx, w, r.erroredData)
	testCases.assert(t.T, w.Body)
}

func (r *RenderSignupTest) HasSignupForm(t *testgroup.T) {
	w := httptest.NewRecorder()
	render.View.Signup(anonCtx, w, render.SignupData{})
	tokens := []string{"Signup", "<form hx-post=\"/signup\"", "</form>"}
	utils.AssertContainsTokens(t.T, w.Result(), tokens)
}

func (r *RenderSignupTest) HasProperLayoutMarkers(t *testgroup.T) {
	layoutTestCases.assert(t.T, func(r render.Render) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			r.Signup(anonCtx, w, render.SignupData{})
		}
	})
}

func (r *RenderSignupTest) PreGroup(t *testgroup.T) {
	form := contracts.SignupForm{
		Username:     contracts.Field{Errors: errs},
		Password:     contracts.Field{Errors: errs},
		Confirmation: contracts.Field{Errors: errs},
	}

	r.erroredData = render.SignupData{form}
}

func TestRenderSignup(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(RenderSignupTest))
}

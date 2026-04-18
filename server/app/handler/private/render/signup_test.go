package render_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type RenderSignupTest struct {
	errors      contracts.Errors
	erroredData render.SignupData
}

func (r *RenderSignupTest) RendersValidationErrors(t *testgroup.T) {
	testCases := fieldErrorsTestCases{{
		"section#username ul.errors",
		r.errors,
	}, {

		"section#password ul.errors",
		r.errors,
	}, {

		"section#confirmation ul.errors",
		r.errors,
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
	errs := contracts.Errors{errors.New("kek-1"), errors.New("kek-2")}
	form := contracts.SignupForm{
		Username:     contracts.Field{Errors: errs},
		Password:     contracts.Field{Errors: errs},
		Confirmation: contracts.Field{Errors: errs},
	}

	r.errors = errs
	r.erroredData = render.SignupData{form}
}

func TestRenderSignup(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(RenderSignupTest))
}

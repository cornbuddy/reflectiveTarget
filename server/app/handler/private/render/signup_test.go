package render_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
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
	type testCase struct {
		selector string
	}

	testCases := []testCase{{
		selector: "section#username ul.errors",
	}, {

		selector: "section#password ul.errors",
	}, {

		selector: "section#confirmation ul.errors",
	}}

	w := httptest.NewRecorder()
	render.View.Signup(userCtx, w, r.erroredData)
	doc, err := goquery.NewDocumentFromReader(w.Body)
	t.Require.NoError(err)

	for _, tc := range testCases {
		ul := doc.Find(tc.selector)
		t.NotEmpty(ul.Nodes, tc.selector)

		ul.Children().Each(func(_ int, li *goquery.Selection) {
			err := errors.New(li.Text())
			t.Contains(r.errors, err, tc.selector)
		})
	}
}

func (r *RenderSignupTest) HasSignupForm(t *testgroup.T) {
	w := httptest.NewRecorder()
	render.View.Signup(anonCtx, w, render.SignupData{})
	tokens := []string{"Signup", "<form hx-post=\"/signup\"", "</form>"}
	utils.AssertContainsTokens(t.T, w.Result(), tokens)
}

func (r *RenderSignupTest) HasProperLayoutMarkers(t *testgroup.T) {
	testCases.run(t.T, func(r render.Render) func(w http.ResponseWriter) {
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

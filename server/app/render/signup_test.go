package render_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
	"github.com/cornbuddy/reflectiveTarget/server/app/render"
)

type RenderSignupTest struct {
	errors      forms.Errors
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
	assertContainsTokens(t.T, w.Body, tokens)
}

func (r *RenderSignupTest) HasProperLayoutMarkers(t *testgroup.T) {
	testCases.run(t.T, func(r render.Render) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			r.Signup(anonCtx, w, render.SignupData{})
		}
	})
}

func (r *RenderSignupTest) PreGroup(t *testgroup.T) {
	errs := forms.Errors{errors.New("kek-1"), errors.New("kek-2")}
	form := forms.SignupForm{
		Username:     forms.Field{Errors: errs},
		Password:     forms.Field{Errors: errs},
		Confirmation: forms.Field{Errors: errs},
	}

	r.errors = errs
	r.erroredData = render.SignupData{form}
}

func TestRenderSignup(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(RenderSignupTest))
}

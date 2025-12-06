package render_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
	"github.com/cornbuddy/reflectiveTarget/server/app/render"
)

func TestSignupFormValidationErrorsShouldBeRendered(t *testing.T) {
	t.Parallel()

	type testCase struct {
		selector string
	}

	testCases := []testCase{{
		selector: "form ul#username-errors",
	}, {

		selector: "form ul#password-errors",
	}, {

		selector: "form ul#confirmation-errors",
	}}

	ctx := context.TODO()
	w := httptest.NewRecorder()
	errs := forms.Errors{errors.New("kek")}
	form := forms.SignupForm{
		Username:     forms.Field{Errors: errs},
		Password:     forms.Field{Errors: errs},
		Confirmation: forms.Field{Errors: errs},
	}

	render.View.Signup(ctx, w, form)
	doc, err := goquery.NewDocumentFromReader(w.Body)
	require.NoError(t, err)

	html, err := doc.Html()
	require.NoError(t, err)
	t.Logf("kekeke: %s", html)

	for _, tc := range testCases {
		ul := doc.Find(tc.selector)
		assert.NotEmpty(t, ul.Nodes, tc.selector)

		ul.Children().Each(func(_ int, li *goquery.Selection) {
			err := errors.New(li.Text())
			assert.Contains(t, errs, err, tc.selector)
		})
	}

}

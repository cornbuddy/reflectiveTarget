package render_test

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"
)

func TestSignupFormValidationErrorsShouldBeRendered(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	w := httptest.NewRecorder()
	errs := []error{errors.New("kek")}
	render.View.Signup(ctx, w, errs)

	doc, err := goquery.NewDocumentFromReader(w.Body)
	require.NoError(t, err)

	ul := doc.Find("form ul#username-errors")
	require.NotNil(t, ul)
	require.NotEmpty(t, ul.Nodes)

	ul.Children().Each(func(_ int, s *goquery.Selection) {
		err := errors.New(s.Text())
		assert.NotContains(t, errs, err)
	})
}

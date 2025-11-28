package render_test

import (
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLayoutRenderer(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	w := httptest.NewRecorder()
	render.Layout.Index(ctx, w, nil)
	raw, err := io.ReadAll(w.Body)
	require.NoError(t, err)

	body := string(raw)
	assert.Contains(t, body, "html")
	assert.Contains(t, body, "h1")
}

package handlers

import (
	"io"
	"net/http"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldRenderIndexPage(t *testing.T) {
	t.Parallel()

	get := http.MethodGet
	resp := utils.MakeRequest("", get, "/", indexRouter.Get, nil)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	t.Cleanup(func() { resp.Body.Close() })

	assert.Contains(t, string(body), "Reflective target")
}

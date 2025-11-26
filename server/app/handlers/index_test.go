package handlers

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestShouldRenderIndexPage(t *testing.T) {
	t.Parallel()

	resp := utils.MakeRequest("", http.MethodGet, "/", router, nil)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	t.Cleanup(func() { resp.Body.Close() })
	body := string(bytes)

	for _, tc := range []struct {
		msg  string
		want string
	}{{
		msg:  "should contain title",
		want: "Reflective target",
	}, {
		msg:  "should contain navigation",
		want: "<nav>",
	}, {
		msg:  "should contain sidebar toggler",
		want: "id=\"sidebar-toggler\"",
	}} {
		assert.Contains(t, body, tc.want, tc.msg)
	}
}

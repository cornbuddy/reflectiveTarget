package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShotsShouldBeValidated(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc string
		body ShotsRequest
	}

	testCases := []testCase{{
		desc: "should fail when coordinates are greater than 100",
		body: ShotsRequest{101, 101},
	}, {
		desc: "should fail when coordinates are less than 0",
		body: ShotsRequest{-1, -1},
	}}

	for _, tc := range testCases {
		var body bytes.Buffer
		require.NoError(t, json.NewEncoder(&body).Encode(tc.body))

		ct := "application/json"
		handle := shotsRouter.Post
		method := http.MethodPost
		resp := utils.MakeRequest(ct, method, "/", handle, &body)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, tc.desc)
	}
}

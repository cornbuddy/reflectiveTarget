package handlers

import (
	"bytes"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/private/model"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShotsShouldBeSavedIfValid(t *testing.T) {
	t.Parallel()

	shot := model.Shot{
		X: rand.IntN(101),
		Y: rand.IntN(101),
	}

	var body bytes.Buffer
	require.NoError(t, json.NewEncoder(&body).Encode(
		model.ShotsRequest{Shots: []model.Shot{shot}},
	))

	ct := "application/json"
	handle := shots.Post
	method := http.MethodPost
	resp := utils.MakeRequest(ct, method, "/", handle, &body)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	res, err := db.Query(
		"SELECT * FROM shots WHERE x = $1 AND y = $2",
		shot.X, shot.Y,
	)
	require.NoError(t, err)
	assert.True(t, res.Next(), "should save shots")
}

func TestShotsShouldBeValidated(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc string
		body model.ShotsRequest
	}

	testCases := []testCase{{
		desc: "should fail when coordinates are greater than 100",
		body: model.ShotsRequest{Shots: []model.Shot{{X: 101, Y: 101}}},
	}, {
		desc: "should fail when coordinates are less than 0",
		body: model.ShotsRequest{Shots: []model.Shot{{X: -1, Y: -1}}},
	}}

	for _, tc := range testCases {
		var body bytes.Buffer
		require.NoError(t, json.NewEncoder(&body).Encode(tc.body))

		ct := "application/json"
		handle := shots.Post
		method := http.MethodPost
		resp := utils.MakeRequest(ct, method, "/", handle, &body)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, tc.desc)
	}
}

func TestShotsShouldFailIfRequestIsMalformed(t *testing.T) {
	t.Parallel()

	var body bytes.Buffer
	_, err := body.Write([]byte("kek"))
	require.NoError(t, err)

	ct := "application/json"
	handle := shots.Post
	method := http.MethodPost
	resp := utils.MakeRequest(ct, method, "/", handle, &body)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

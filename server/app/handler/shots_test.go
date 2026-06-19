package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

const shotsUrl = "/api/target/1/shots"

func TestShouldReturnNoShotsForEmptyTarget(t *testing.T) {
	t.Parallel()

	user, err := makeTestUser(db)
	require.NoError(t, err)

	targetID, err := makeTestTarget(db, int(user.ID))
	require.NoError(t, err)

	ct := "application/json"
	url := fmt.Sprintf("/api/target/%v/shots", targetID)
	resp, body, err := utils.MakeRequest(ct, http.MethodGet, url, router, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, ctAppJson, resp.Header.Get("Content-Type"))

	var shots contracts.ShotsResponse
	require.NoError(t, json.Unmarshal([]byte(body), &shots))
	assert.Empty(t, shots.Shots)
}

func TestShotsShould404TargetDoesNotExist(t *testing.T) {
	t.Parallel()

	ct := "application/json"
	url := "/api/target/69/shots"
	resp, _, err := utils.MakeRequest(ct, http.MethodGet, url, router, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestShotsShouldBeSavedIfValid(t *testing.T) {
	t.Parallel()

	user, err := makeTestUser(db)
	require.NoError(t, err)

	targetID, err := makeTestTarget(db, int(user.ID))
	require.NoError(t, err)

	shot := valueobjects.Shot{
		X: rand.IntN(101),
		Y: rand.IntN(101),
	}

	var shots bytes.Buffer
	require.NoError(t, json.NewEncoder(&shots).Encode(
		contracts.ShotsRequest{Shots: []valueobjects.Shot{shot}},
	))

	url := fmt.Sprintf("/api/target/%v/shots", targetID)
	resp, body, err := utils.MakeRequestWithCookies(
		"application/json", http.MethodPost, url, router, &shots,
		&http.Cookie{Name: constants.SessionCookieName, Value: "kek"},
	)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.JSONEq(t, ctAppJson, resp.Header.Get("Content-Type"))
	assert.Equal(t, "ok", body, "should save shots")

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
		body contracts.ShotsRequest
	}

	testCases := []testCase{{
		desc: "should fail when coordinates are greater than 100",
		body: contracts.ShotsRequest{
			Shots: []valueobjects.Shot{{X: 101, Y: 101}},
		},
	}, {
		desc: "should fail when coordinates are less than 0",
		body: contracts.ShotsRequest{
			Shots: []valueobjects.Shot{{X: -1, Y: -1}},
		},
	}}

	for _, tc := range testCases {
		var body bytes.Buffer
		require.NoError(t, json.NewEncoder(&body).Encode(tc.body))

		ct := "application/json"
		method := http.MethodPost
		resp, _, err := utils.MakeRequest(ct, method, shotsUrl, router, &body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode, tc.desc)
	}
}

func TestShotsShouldFailIfRequestIsMalformed(t *testing.T) {
	t.Parallel()

	var body bytes.Buffer
	_, err := body.Write([]byte("kek"))
	require.NoError(t, err)

	ct := "application/json"
	method := http.MethodPost
	resp, _, err := utils.MakeRequest(ct, method, shotsUrl, router, &body)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

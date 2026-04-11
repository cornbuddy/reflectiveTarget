package handler

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/bloomberg/go-testgroup"

	appconst "github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type TargetsSuite struct {
	session        []*http.Cookie
	ownedTargets   aggregations.Targets
	foreignTargets aggregations.Targets
}

func (s *TargetsSuite) ShouldRejectInvalidTarget(t *testgroup.T) {
	t.Fail("not implemented")
	// TODO: make post request with invalid target as form
	// TODO: ensure that POST /targets/new responds with 400 code
}

func (s *TargetsSuite) ShouldAddTargetIfValid(t *testgroup.T) {
	t.Fail("not implemented")
	// TODO: make post request with valid target as form
	// TODO: ensure that /targets endpoint contains new target name
}

func (s *TargetsSuite) ShouldRespondOnValidCreds(t *testgroup.T) {
	r := utils.MakeRequestWithCookies(
		"", http.MethodGet, "/targets/new", router, nil, s.session...,
	)
	t.Equal(http.StatusOK, r.StatusCode)
}

func (s *TargetsSuite) ShouldListTargetsForOwner(t *testgroup.T) {
	r := utils.MakeRequestWithCookies(
		"", http.MethodGet, "/targets", router, nil, s.session...,
	)
	t.Equal(http.StatusOK, r.StatusCode)

	for _, ot := range s.ownedTargets {
		tokens := []string{strconv.Itoa(int(ot.ID)), ot.Name}
		utils.AssertContainsTokens(t.T, r.Body, tokens)
	}

	for _, ft := range s.foreignTargets {
		tokens := []string{strconv.Itoa(int(ft.ID)), ft.Name}
		utils.AssertNotContainsTokens(t.T, r.Body, tokens)
	}
}

func (s *TargetsSuite) PreGroup(t *testgroup.T) {
	owner, err := entities.NewUser("user1", "password")
	t.Require.NoError(err)

	foreigner, err := entities.NewUser("user2", "password")
	t.Require.NoError(err)

	t.Require.NoError(utils.InsertUser(db, owner))
	t.Require.NoError(utils.InsertUser(db, foreigner))

	questions := vo.Questions{
		{Text: "kek1?"},
		{Text: "kek2?"},
	}
	shots := vo.Shots{
		{X: 1, Y: 100},
		{X: 100, Y: 1},
	}
	ownedTargets := aggregations.Targets{aggregations.Target{
		Name:      "test1",
		Owner:     *owner,
		Questions: append(vo.Questions{}, questions...),
		Shots:     append(vo.Shots{}, shots...),
	}, {
		Name:      "test2",
		Owner:     *owner,
		Questions: append(vo.Questions{}, questions...),
		Shots:     append(vo.Shots{}, shots...),
	}}
	t.Require.NoError(utils.InsertTargets(db, ownedTargets))

	foreignTargets := aggregations.Targets{aggregations.Target{
		Name:      "test3",
		Owner:     *foreigner,
		Questions: append(vo.Questions{}, questions...),
		Shots:     append(vo.Shots{}, shots...),
	}}
	t.Require.NoError(utils.InsertTargets(db, foreignTargets))

	token := "kekeke"
	session := sessiondata.SessionData{
		IsAuthenticated: true,
		UserID:          owner.ID,
		Username:        owner.Username,
	}
	t.Require.NoError(sessionStore.Update(ctx, token, session))

	s.foreignTargets = foreignTargets
	s.ownedTargets = ownedTargets
	s.session = []*http.Cookie{{
		Name:  appconst.SessionCookieName,
		Value: token,
	}}
}

func TestTargetsHandler(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetsSuite))
}

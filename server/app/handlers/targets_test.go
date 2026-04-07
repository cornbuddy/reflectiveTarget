package handlers

import (
	"io"
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

func (s *TargetsSuite) ShouldContainHtmlForm(t *testgroup.T) {
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

	data, err := io.ReadAll(r.Body)
	t.Require.NoError(err)

	t.Cleanup(func() {
		t.Require.NoError(r.Body.Close())
	})

	body := string(data)
	for _, ot := range s.ownedTargets {
		t.Contains(body, strconv.Itoa(int(ot.ID)))
		t.Contains(body, ot.Name)
	}

	for _, ft := range s.foreignTargets {
		t.NotContains(body, strconv.Itoa(int(ft.ID)))
		t.NotContains(body, ft.Name)
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

package handler

import (
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/bloomberg/go-testgroup"

	appconst "github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type TargetsSuite struct {
	handler        http.HandlerFunc
	session        []*http.Cookie
	ownedTargets   aggregations.Targets
	foreignTargets aggregations.Targets
}

func (s *TargetsSuite) ShouldRenderFormWithTarget(t *testgroup.T) {
	target := s.ownedTargets[0]
	url := fmt.Sprintf("/targets/%d", target.ID)
	t.HTTPStatusCode(s.handler, http.MethodGet, url, nil, http.StatusOK)

	r, body, err := utils.MakeRequest("", http.MethodGet, url, s.handler, nil)
}

func (s *TargetsSuite) ShouldRejectInvalidTarget(t *testgroup.T) {
	type testCase struct {
		desc     string
		form     neturl.Values
		contains string
	}

	testCases := []testCase{{
		"rejects target with existing name name",
		neturl.Values{
			"name":       []string{s.ownedTargets[0].Name},
			"question_0": []string{"q1?"},
			"question_1": []string{"q2?"},
		},
		validators.ErrTargetAlreadyExists.Error(),
	}, {
		"rejects target with no questions",
		neturl.Values{
			"name": []string{"no questions"},
		},
		validators.ErrEmpty.Error(),
	}, {
		"rejects target with repeated questions",
		neturl.Values{
			"name":       []string{"same question twice"},
			"question_0": []string{"q1?"},
			"question_1": []string{"q1?"},
		},
		validators.ErrRepeatedQuestion.Error(),
	}}

	status := http.StatusBadRequest
	method := http.MethodPost
	url := "/targets/new"
	ct := "application/x-www-form-urlencoded"
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testgroup.T) {
			t.Parallel()

			form := strings.NewReader(tc.form.Encode())
			r, body, err := utils.MakeRequest(ct, method, url, s.handler, form)
			t.Require.NoError(err)
			t.Equal(status, r.StatusCode)
			t.Contains(body, tc.contains)
		})
	}
}

func (s *TargetsSuite) ShouldAddTargetIfValid(t *testgroup.T) {
	name := "valid target"
	form := neturl.Values{
		"name":       []string{name},
		"question_0": []string{"q1?"},
		"question_1": []string{"q2?"},
	}
	t.HTTPStatusCode(
		s.handler, http.MethodPost, "/targets/new", form, http.StatusSeeOther,
	)
	t.HTTPBodyContains(s.handler, http.MethodGet, "/targets", nil, name)
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
		utils.AssertContainsTokens(t.T, r, tokens)

	}

	for _, ft := range s.foreignTargets {
		tokens := []string{strconv.Itoa(int(ft.ID)), ft.Name}
		utils.AssertNotContainsTokens(t.T, r, tokens)
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
	s.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, c := range s.session {
			r.AddCookie(c)
		}

		router(w, r)
	})
}

func TestTargetsHandler(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetsSuite))
}

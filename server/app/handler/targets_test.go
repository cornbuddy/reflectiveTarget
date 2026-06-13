package handler

import (
	"fmt"
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
	owner          *entities.User
	ownedTargets   aggregations.Targets
	foreignTargets aggregations.Targets
}

func (s *TargetsSuite) UpdateShouldBeIdempotent(t *testgroup.T) {
	question := vo.Question{Text: utils.MakeRandomString(5)}
	target := aggregations.Target{
		Name:      utils.MakeRandomString(5),
		Owner:     *s.owner,
		Questions: vo.Questions{question},
	}
	t.Require.NoError(utils.InsertTarget(db, &target))

	url := fmt.Sprintf("/targets/%d", target.ID)
	form := strings.NewReader(neturl.Values{
		"name":             []string{target.Name},
		"question_0_id":    []string{strconv.Itoa(int(question.ID))},
		"question_0_value": []string{question.Text},
	}.Encode())

	ct := "application/x-www-form-urlencoded"
	r, body, err := utils.MakeRequest(ct, http.MethodPut, url, s.handler, form)
	t.Require.NoError(err)
	t.Equal(http.StatusSeeOther, r.StatusCode)
	t.Equal("target updated", body)
}

func (s *TargetsSuite) ShouldUpdateExistingTarget(t *testgroup.T) {
	oldQstn := "single question"
	question := vo.Question{Text: oldQstn}
	target := aggregations.Target{
		Name:      utils.MakeRandomString(5),
		Owner:     *s.owner,
		Questions: vo.Questions{question},
	}
	t.Require.NoError(utils.InsertTarget(db, &target))

	newQstn := "updated question, still single"
	url := fmt.Sprintf("/targets/%d", target.ID)
	form := strings.NewReader(neturl.Values{
		"name":             []string{target.Name},
		"question_0_id":    []string{strconv.Itoa(int(question.ID))},
		"question_0_value": []string{newQstn},
	}.Encode())

	ct := "application/x-www-form-urlencoded"
	r, body, err := utils.MakeRequest(ct, http.MethodPut, url, s.handler, form)
	t.Require.NoError(err)
	t.Equal(http.StatusSeeOther, r.StatusCode)
	t.Equal("target updated", body)

	r, body, err = utils.MakeRequest(ct, http.MethodGet, url, s.handler, form)
	t.Require.NoError(err)
	t.Contains(body, target.Name)
	t.Contains(body, newQstn)
	t.NotContains(body, oldQstn, "should update question, not create a new one")
}

func (s *TargetsSuite) ShouldRenderFormWithTarget(t *testgroup.T) {
	target := s.ownedTargets[0]
	url := fmt.Sprintf("/targets/%d", target.ID)
	r, _, err := utils.MakeRequest("", http.MethodGet, url, s.handler, nil)
	t.Require.NoError(err)
	t.Equal(http.StatusOK, r.StatusCode)

	tokens := []string{target.Name, "Update"}
	for _, q := range target.Questions {
		tokens = append(tokens, q.Text)
	}

	utils.AssertContainsTokens(t.T, r, tokens)
	utils.AssertNotContainsTokens(t.T, r, []string{"Create", "Add question"})
}

func (s *TargetsSuite) ShouldRespondWithNotFoundIfNoTarget(t *testgroup.T) {
	url := "/targets/69"
	r, body, err := utils.MakeRequest("", http.MethodGet, url, s.handler, nil)
	t.Require.NoError(err)
	t.Equal(http.StatusNotFound, r.StatusCode)
	t.Contains(body, "Not Found")
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
			"name":             []string{s.ownedTargets[0].Name},
			"question_0_value": []string{"q1?"},
			"question_1_value": []string{"q2?"},
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
			"name":             []string{"same question twice"},
			"question_0_value": []string{"q1?"},
			"question_1_value": []string{"q1?"},
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
		"name":             []string{name},
		"question_0_value": []string{"q1?"},
		"question_1_value": []string{"q2?"},
	}
	t.HTTPStatusCode(
		s.handler, http.MethodPost, "/targets/new", form, http.StatusSeeOther,
	)
	t.HTTPBodyContains(s.handler, http.MethodGet, "/targets", nil, name)
}

func (s *TargetsSuite) ShouldRespondOnValidCreds(t *testgroup.T) {
	r, _, err := utils.MakeRequest(
		"", http.MethodGet, "/targets/new", s.handler, nil,
	)
	t.Require.NoError(err)
	t.Equal(http.StatusOK, r.StatusCode)
}

func (s *TargetsSuite) ShouldListTargetsForOwner(t *testgroup.T) {
	r, _, err := utils.MakeRequest(
		"", http.MethodGet, "/targets", s.handler, nil,
	)
	t.Require.NoError(err)
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

	s.owner = owner
	s.ownedTargets = ownedTargets
	s.foreignTargets = foreignTargets
	s.handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookies := []*http.Cookie{{
			Name:  appconst.SessionCookieName,
			Value: token,
		}}

		for _, c := range cookies {
			r.AddCookie(c)
		}

		router(w, r)
	})
}

func TestTargetsHandler(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetsSuite))
}

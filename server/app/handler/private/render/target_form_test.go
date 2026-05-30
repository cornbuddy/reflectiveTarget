package render_test

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type TargetFormTest struct{}

func (u *TargetFormTest) RendersErrorsProperly(t *testgroup.T) {
	t.Skip()

	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	data := render.TargetFormData{
		TargetForm: contracts.TargetForm{
			Name: contracts.Field{
				Value:  "",
				Errors: contracts.Errors{err1},
			},
			Questions: contracts.Fields{{
				Value:  "kek1",
				Errors: contracts.Errors{err2},
			}, {
				Value:  "kek12",
				Errors: contracts.Errors{err1, err2},
			}, {
				Value: "kek0",
			}},
		},
	}

	testCases := fieldErrorsTestCases{{
		"section#name ul.errors",
		contracts.Errors{err1},
	}, {
		"section#question_0 ul.errors",
		contracts.Errors{err2},
	}, {
		"section#question_1 ul.errors",
		contracts.Errors{err1, err2},
	}, {
		"section#question_2 ul.errors",
		contracts.Errors{},
	}}

	w := httptest.NewRecorder()
	render.View.TargetForm(anonCtx, w, data)
	testCases.assert(t.T, w.Body)
}

func (u *TargetFormTest) RendersProperViewForm(t *testgroup.T) {
	type testCase struct {
		desc   string
		data   render.TargetFormData
		tokens []string
		assert func(*testing.T, *http.Response, []string)
	}

	addQuestion := "<button hx-on:click=\"addQuestion()\">Add question</button>"
	target := aggregations.Target{
		Name: "kek",
		ID:   69,
		Questions: valueobjects.Questions{
			{Text: "fails there"}, {Text: "kek2?"},
		},
	}
	newQuestions := make(contracts.Fields, len(target.Questions))
	existingQuestions := make(contracts.Fields, len(target.Questions))
	for i, q := range target.Questions {
		newQuestions[i] = contracts.Field{Value: q.Text}
		existingQuestions[i] = contracts.Field{
			ID:    valueobjects.ID(rand.IntN(100)),
			Value: q.Text,
		}
	}
	newForm := contracts.TargetForm{
		Name:      contracts.Field{Value: target.Name},
		Questions: newQuestions,
	}
	existingForm := contracts.TargetForm{
		Name:      contracts.Field{Value: target.Name},
		Questions: existingQuestions,
	}
	testCases := []testCase{{
		"saved questions' inputs have data-id property",
		render.TargetFormData{existingForm, target.ID},
		[]string{
			fmt.Sprintf("<h2>%s</h2>", target.Name),
			"<form",
			fmt.Sprintf("hx-put=\"/targets/%d\"", target.ID),
			"hx-trigger=\"submit\"",
			"hx-target=\"main\"",
			"<button type=\"submit\">Update</button>",
			fmt.Sprintf(
				`<input name="name" placeholder="Target name" value="%s"`,
				target.Name,
			),
			fmt.Sprintf(
				`<input name="%s" placeholder="Question 1" data-id="%d" value="%s">`,
				"question_0",
				existingQuestions[0].ID,
				existingQuestions[0].Value,
			),
			fmt.Sprintf(
				`<input name="%s" placeholder="Question 2" data-id="%d" value="%s">`,
				"question_1",
				existingQuestions[1].ID,
				existingQuestions[1].Value,
			),
		},
		utils.AssertContainsTokens,
	}, {
		"empty target",
		render.TargetFormData{},
		[]string{
			"New target",
			"<form",
			"hx-post=\"/targets/new\"",
			"hx-trigger=\"submit\"",
			"hx-target=\"main\"",
			"<button type=\"submit\">Create</button>",
			"<input name=\"name\" placeholder=\"Target name\" value=\"\"",
			"<input name=\"question_0\" placeholder=\"Question 1\" data-id=\"0\">",
			addQuestion,
		},
		utils.AssertContainsTokens,
	}, {
		"non empty target contains",
		render.TargetFormData{newForm, target.ID},
		[]string{
			fmt.Sprintf("<h2>%s</h2>", target.Name),
			"<form",
			fmt.Sprintf("hx-put=\"/targets/%d\"", target.ID),
			"hx-trigger=\"submit\"",
			"hx-target=\"main\"",
			"<button type=\"submit\">Update</button>",
			fmt.Sprintf(
				"<input name=\"name\" placeholder=\"Target name\" value=\"%s\"",
				target.Name,
			),
			fmt.Sprintf(
				"<input name=\"%s\" placeholder=\"Question 1\" data-id=\"0\" value=\"%s\">",
				"question_0",
				target.Questions[0].Text,
			),
			fmt.Sprintf(
				"<input name=\"%s\" placeholder=\"Question 2\" data-id=\"0\" value=\"%s\">",
				"question_1",
				target.Questions[1].Text,
			),
		},
		utils.AssertContainsTokens,
	}, {
		"non empty target doesn't contain",
		render.TargetFormData{newForm, target.ID},
		[]string{addQuestion},
		utils.AssertNotContainsTokens,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testgroup.T) {
			w := httptest.NewRecorder()
			render.View.TargetForm(anonCtx, w, tc.data)
			tc.assert(t.T, w.Result(), tc.tokens)
		})
	}
}

func (u *TargetFormTest) HasProperLayoutMarkers(t *testgroup.T) {
	layoutTestCases.assert(t.T, func(r render.Render) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			r.TargetForm(anonCtx, w, render.TargetFormData{})
		}
	})
}

func TestUpdateTarget(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetFormTest))
}

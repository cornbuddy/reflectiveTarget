package render_test

import (
	"errors"
	"fmt"
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
	formQuestions := make(contracts.Fields, len(target.Questions))
	for i, q := range target.Questions {
		formQuestions[i] = contracts.Field{Value: q.Text}
	}
	form := contracts.TargetForm{
		Name:      contracts.Field{Value: target.Name},
		Questions: formQuestions,
	}
	testCases := []testCase{{
		"empty target",
		render.TargetFormData{},
		[]string{
			"New target",
			"<form",
			"hx-post=\"/target/new\"",
			"hx-trigger=\"submit\"",
			"hx-target=\"main\"",
			"<button type=\"submit\">Create</button>",
			"<input name=\"name\" placeholder=\"Target name\" value=\"\"",
			"<input name=\"question_0\" placeholder=\"Question 1\">",
			addQuestion,
		},
		utils.AssertContainsTokens,
	}, {
		"non empty target contains",
		render.TargetFormData{form, target.ID},
		[]string{
			fmt.Sprintf("<h2>%s</h2>", target.Name),
			"<form",
			fmt.Sprintf("hx-put=\"/target/%d\"", target.ID),
			"hx-trigger=\"submit\"",
			"hx-target=\"main\"",
			"<button type=\"submit\">Update</button>",
			fmt.Sprintf(
				"<input name=\"name\" placeholder=\"Target name\" value=\"%s\"",
				target.Name,
			),
			fmt.Sprintf(
				"<input name=\"question_0\" placeholder=\"Question 1\" value=\"%s\">",
				target.Questions[0].Text,
			),
			fmt.Sprintf(
				"<input name=\"question_1\" placeholder=\"Question 2\" value=\"%s\">",
				target.Questions[1].Text,
			),
		},
		utils.AssertContainsTokens,
	}, {
		"non empty target doesn't contain",
		render.TargetFormData{form, target.ID},
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

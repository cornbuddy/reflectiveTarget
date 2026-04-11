package render_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/render"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type UpdateTargetTest struct{}

func (u *UpdateTargetTest) RendersProperViewFor(t *testgroup.T) {
	type testCase struct {
		desc   string
		data   render.TargetData
		tokens []string
		assert func(*testing.T, io.Reader, []string)
	}

	addQuestion := "<button hx-on:click=\"addQuestion()\">Add question</button>"
	target := aggregations.Target{
		Name: "kek",
		ID:   69,
		Questions: valueobjects.Questions{
			{Text: "kek1?"}, {Text: "kek2?"},
		},
	}
	testCases := []testCase{{
		"empty target",
		render.TargetData{},
		[]string{
			"New target",
			"<form",
			"hx-post=\"/target/new\"",
			"hx-trigger=\"submit\"",
			"hx-target=\"main\"",
			"<button type=\"submit\">Create</button>",
			"<input name=\"name\" placeholder=\"Target name\" value=\"\">",
			"<input name=\"question_0\" placeholder=\"Question 1\">",
			addQuestion,
		},
		utils.AssertContainsTokens,
	}, {
		"non empty target contains",
		render.TargetData{target},
		[]string{
			fmt.Sprintf("<h2>%s</h2>", target.Name),
			"<form",
			fmt.Sprintf("hx-put=\"/target/%d\"", target.ID),
			"hx-trigger=\"submit\"",
			"hx-target=\"main\"",
			"<button type=\"submit\">Update</button>",
			fmt.Sprintf(
				"<input name=\"name\" placeholder=\"Target name\" value=\"%s\">",
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
		render.TargetData{target},
		[]string{addQuestion},
		utils.AssertNotContainsTokens,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testgroup.T) {
			w := httptest.NewRecorder()
			render.View.UpdateTarget(anonCtx, w, tc.data)
			tc.assert(t.T, w.Body, tc.tokens)
		})
	}
}

func (u *UpdateTargetTest) HasProperLayoutMarkers(t *testgroup.T) {
	testCases.run(t.T, func(r render.Render) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			r.UpdateTarget(anonCtx, w, render.TargetData{})
		}
	})
}

func TestUpdateTarget(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(UpdateTargetTest))
}

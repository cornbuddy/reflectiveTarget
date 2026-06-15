package valueobjects_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

func TestQuestionStringer(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc     string
		question vo.Question
		want     string
	}

	testCases := []testCase{{
		"default",
		vo.Question{},
		"{ID: 0, Text: ''}",
	}, {
		"not default",
		vo.Question{ID: 69, Text: "kek"},
		"{ID: 69, Text: 'kek'}",
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.question.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestQuestionsStringer(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc      string
		questions vo.Questions
		want      string
	}

	testCasess := []testCase{{
		"empty",
		vo.Questions{},
		"[]",
	}, {
		"not empty",
		vo.Questions{vo.Question{"1", 1}, vo.Question{"2", 2}},
		"[{ID: 1, Text: '1'}, {ID: 2, Text: '2'}]",
	}}

	for _, tc := range testCasess {
		t.Run(tc.desc, func(t *testing.T) {
			got := tc.questions.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

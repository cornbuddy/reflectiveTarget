package contracts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

func TestQuestionFieldMethods(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc      string
		field     contracts.QuestionField
		wantValue string
		wantID    string
	}

	const (
		qPref    = contracts.QuestionFieldPrefix
		qValPost = contracts.QuestionFieldValuePostfix
		qIDPost  = contracts.QuestionFieldIDPostfix
	)

	testCases := []testCase{{
		"uses 0 by default",
		contracts.QuestionField{},
		qPref + "_0_" + qValPost,
		qPref + "_0_" + qIDPost,
	}, {
		"uses id in the middle",
		contracts.QuestionField{contracts.Field{ID: 69}},
		qPref + "_69_" + qValPost,
		qPref + "_69_" + qIDPost,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.wantValue, tc.field.NameValue())
			assert.Equal(t, tc.wantID, tc.field.NameID())
		})
	}
}

package contracts_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

var (
	_q        = contracts.QuestionField{}
	formQ0ID  = _q.NameID(0)
	formQ0Val = _q.NameValue(0)
	formQ1ID  = _q.NameID(1)
	formQ1Val = _q.NameValue(1)
)

func TestTargetFormStringer(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc   string
		form   contracts.TargetForm
		regexp string
	}

	testCases := []testCase{{
		"defaults",
		contracts.TargetForm{},
		`^Name: {.+}, Questions: \[\]$`,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			got := tc.form.String()
			assert.Regexp(t, tc.regexp, got)
		})
	}
}

func TestNewTargetFormFromTarget(t *testing.T) {
	t.Parallel()

	type testCase struct {
		target aggregations.Target
		form   contracts.TargetForm
	}

	testCases := []testCase{{
		aggregations.Target{
			Name: kek,
			Questions: valueobjects.Questions{{
				ID:   1,
				Text: kek1,
			}, {
				ID:   2,
				Text: kek2,
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: kek},
			Questions: contracts.QuestionFields{{
				contracts.Field{
					ID:    1,
					Value: kek1,
				},
			}, {
				contracts.Field{
					ID:    2,
					Value: kek2,
				},
			}},
		},
	}}

	for _, tc := range testCases {
		got := contracts.NewTargetFormFromTarget(tc.target)
		assert.Equal(t, tc.form, got)
	}
}

func TestNewTargetForm(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc        string
		form        url.Values
		want        *contracts.TargetForm
		err         error
		errContains []string
	}

	testCases := []testCase{{
		"should handle empty form",
		url.Values{},
		&contracts.TargetForm{},
		nil,
		nil,
	}, {
		"should handle name only",
		url.Values{formName: []string{kek}},
		&contracts.TargetForm{
			Name: contracts.Field{Value: kek},
		},
		nil,
		nil,
	}, {
		"should handle question value only",
		url.Values{formQ0Val: []string{kek}},
		&contracts.TargetForm{
			Questions: contracts.QuestionFields{{
				contracts.Field{Value: kek}},
			},
		},
		nil,
		nil,
	}, {
		"should error if question id is filled with crap",
		url.Values{
			formName:  []string{qek},
			formQ0ID:  []string{"not a number"},
			formQ0Val: []string{kek1},
		},
		nil,
		contracts.ErrInvalidFieldValue,
		[]string{formQ0ID},
	}, {
		"should handle full form",
		url.Values{
			formName:  []string{qek},
			formQ0Val: []string{kek1},
			formQ1Val: []string{kek2},
		},
		&contracts.TargetForm{
			Name: contracts.Field{Value: qek},
			Questions: contracts.QuestionFields{
				{contracts.Field{Value: kek1}},
				{contracts.Field{Value: kek2}},
			},
		},
		nil,
		nil,
	}, {
		"should handle question id",
		url.Values{
			formName:  []string{qek},
			formQ0ID:  []string{"69"},
			formQ0Val: []string{kek1},
		},
		&contracts.TargetForm{
			Name: contracts.Field{Value: qek},
			Questions: contracts.QuestionFields{{
				contracts.Field{ID: 69, Value: kek1},
			}},
		},
		nil,
		nil,
	}, {
		"should handle multiple questions ids",
		url.Values{
			formName:  []string{qek},
			formQ0ID:  []string{"69"},
			formQ0Val: []string{kek1},
			formQ1ID:  []string{"420"},
			formQ1Val: []string{kek2},
		},
		&contracts.TargetForm{
			Name: contracts.Field{Value: qek},
			Questions: contracts.QuestionFields{
				{contracts.Field{ID: 69, Value: kek1}},
				{contracts.Field{ID: 420, Value: kek2}},
			},
		},
		nil,
		nil,
	}, {
		"should handle new and existing questions",
		url.Values{
			formName:  []string{qek},
			formQ0ID:  []string{"69"},
			formQ0Val: []string{kek1},
			formQ1Val: []string{kek2},
		},
		&contracts.TargetForm{
			Name: contracts.Field{Value: qek},
			Questions: contracts.QuestionFields{
				{contracts.Field{ID: 69, Value: kek1}},
				{contracts.Field{Value: kek2}},
			},
		},
		nil,
		nil,
	}, {
		"should handle a lot of questions",
		url.Values{
			formName:            []string{qek},
			formQ0ID:            []string{"1"},
			formQ0Val:           []string{"1"},
			formQ1ID:            []string{"2"},
			formQ1Val:           []string{"2"},
			"question_2_id":     []string{"3"},
			"question_2_value":  []string{"3"},
			"question_3_id":     []string{"4"},
			"question_3_value":  []string{"4"},
			"question_4_id":     []string{"5"},
			"question_4_value":  []string{"5"},
			"question_5_id":     []string{"6"},
			"question_5_value":  []string{"6"},
			"question_6_id":     []string{"7"},
			"question_6_value":  []string{"7"},
			"question_7_id":     []string{"8"},
			"question_7_value":  []string{"8"},
			"question_8_id":     []string{"9"},
			"question_8_value":  []string{"9"},
			"question_9_id":     []string{"10"},
			"question_9_value":  []string{"10"},
			"question_10_id":    []string{"11"},
			"question_10_value": []string{"11"},
		},
		&contracts.TargetForm{
			Name: contracts.Field{Value: qek},
			Questions: contracts.QuestionFields{
				{contracts.Field{ID: 1, Value: "1"}},
				{contracts.Field{ID: 2, Value: "2"}},
				{contracts.Field{ID: 3, Value: "3"}},
				{contracts.Field{ID: 4, Value: "4"}},
				{contracts.Field{ID: 5, Value: "5"}},
				{contracts.Field{ID: 6, Value: "6"}},
				{contracts.Field{ID: 7, Value: "7"}},
				{contracts.Field{ID: 8, Value: "8"}},
				{contracts.Field{ID: 9, Value: "9"}},
				{contracts.Field{ID: 10, Value: "10"}},
				{contracts.Field{ID: 11, Value: "11"}},
			},
		},
		nil,
		nil,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			got, err := contracts.NewTargetForm(tc.form)
			require.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.want, got)

			for _, msg := range tc.errContains {
				require.ErrorContains(t, err, msg)
			}
		})
	}
}

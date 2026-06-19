package contracts_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
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
			Name: "kek",
			Questions: valueobjects.Questions{{
				ID:   1,
				Text: "kek1",
			}, {
				ID:   2,
				Text: "kek2",
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: "kek"},
			Questions: contracts.Fields{{
				ID:    1,
				Value: "kek1",
			}, {
				ID:    2,
				Value: "kek2",
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
		url.Values{"name": []string{"kek"}},
		&contracts.TargetForm{
			Name: contracts.Field{Value: "kek"},
		},
		nil,
		nil,
	}, {
		"should handle question value only",
		url.Values{"question_0_value": []string{"kek"}},
		&contracts.TargetForm{
			Questions: []contracts.Field{{Value: "kek"}},
		},
		nil,
		nil,
	}, {
		"should error if question id is filled with crap",
		url.Values{
			"name":             []string{"kek?"},
			"question_0_id":    []string{"not a number"},
			"question_0_value": []string{"kek1"},
		},
		nil,
		contracts.ErrInvalidFieldValue,
		[]string{"question_0_id"},
	}, {
		"should handle full form",
		url.Values{
			"name":             []string{"kek?"},
			"question_0_value": []string{"kek1"},
			"question_1_value": []string{"kek2"},
		},
		&contracts.TargetForm{
			Name: contracts.Field{Value: "kek?"},
			Questions: []contracts.Field{
				{Value: "kek1"},
				{Value: "kek2"},
			},
		},
		nil,
		nil,
	}, {
		"should handle question id",
		url.Values{
			"name":             []string{"kek?"},
			"question_0_id":    []string{"69"},
			"question_0_value": []string{"kek1"},
		},
		&contracts.TargetForm{
			Name: contracts.Field{Value: "kek?"},
			Questions: []contracts.Field{
				{ID: 69, Value: "kek1"},
			},
		},
		nil,
		nil,
	}, {
		"should handle multiple questions ids",
		url.Values{
			"name":             []string{"kek?"},
			"question_0_id":    []string{"69"},
			"question_0_value": []string{"kek1"},
			"question_1_id":    []string{"420"},
			"question_1_value": []string{"kek2"},
		},
		&contracts.TargetForm{
			Name: contracts.Field{Value: "kek?"},
			Questions: []contracts.Field{
				{ID: 69, Value: "kek1"},
				{ID: 420, Value: "kek2"},
			},
		},
		nil,
		nil,
	}, {
		"should handle new and existing questions",
		url.Values{
			"name":             []string{"kek?"},
			"question_0_id":    []string{"69"},
			"question_0_value": []string{"kek1"},
			"question_1_value": []string{"kek2"},
		},
		&contracts.TargetForm{
			Name: contracts.Field{Value: "kek?"},
			Questions: []contracts.Field{
				{ID: 69, Value: "kek1"},
				{Value: "kek2"},
			},
		},
		nil,
		nil,
	}, {
		"should handle a lot of questions",
		url.Values{
			"name":              []string{"kek?"},
			"question_0_id":     []string{"1"},
			"question_0_value":  []string{"1"},
			"question_1_id":     []string{"2"},
			"question_1_value":  []string{"2"},
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
			Name: contracts.Field{Value: "kek?"},
			Questions: []contracts.Field{
				{ID: 1, Value: "1"},
				{ID: 2, Value: "2"},
				{ID: 3, Value: "3"},
				{ID: 4, Value: "4"},
				{ID: 5, Value: "5"},
				{ID: 6, Value: "6"},
				{ID: 7, Value: "7"},
				{ID: 8, Value: "8"},
				{ID: 9, Value: "9"},
				{ID: 10, Value: "10"},
				{ID: 11, Value: "11"},
			},
		},
		nil,
		nil,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := contracts.NewTargetForm(tc.form)
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.want, got)
			for _, msg := range tc.errContains {
				assert.ErrorContains(t, err, msg)
			}
		})
	}
}

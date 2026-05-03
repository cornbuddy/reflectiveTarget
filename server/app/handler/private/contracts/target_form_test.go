package contracts_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

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
				Text: "kek1",
			}, {
				Text: "kek2",
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: "kek"},
			Questions: contracts.Fields{{
				Value: "kek1",
			}, {
				Value: "kek2",
			}},
		},
	}}

	for _, tc := range testCases {
		got := contracts.NewTargetFormFromTarget(tc.target)
		assert.EqualValues(t, tc.form, got)
	}
}

func TestNewTargetForm(t *testing.T) {
	t.Parallel()

	type testCase struct {
		httpForm url.Values
		wantForm contracts.TargetForm
	}

	testCases := []testCase{{
		url.Values{},
		contracts.TargetForm{},
	}, {
		url.Values{"name": []string{"kek"}},
		contracts.TargetForm{
			Name: contracts.Field{Value: "kek"},
		},
	}, {
		url.Values{"question_0": []string{"kek"}},
		contracts.TargetForm{
			Questions: []contracts.Field{{Value: "kek"}},
		},
	}, {
		url.Values{
			"question_0": []string{"kek1"},
			"question_1": []string{"kek2"},
		},
		contracts.TargetForm{
			Questions: []contracts.Field{
				{Value: "kek1"},
				{Value: "kek2"},
			},
		},
	}, {
		url.Values{
			"name":       []string{"kek?"},
			"question_0": []string{"kek1"},
			"question_1": []string{"kek2"},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: "kek?"},
			Questions: []contracts.Field{
				{Value: "kek1"},
				{Value: "kek2"},
			},
		},
	}}

	for _, tc := range testCases {
		gotForm := contracts.NewTargetForm(tc.httpForm)
		assert.Equal(t, tc.wantForm, gotForm)
	}
}

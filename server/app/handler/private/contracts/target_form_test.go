package contracts_test

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

func TestTargetForm(t *testing.T) {
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

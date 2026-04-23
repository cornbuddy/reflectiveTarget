package builders_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/builders"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
)

func TestTargetAssembler(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc       string
		form       contracts.TargetForm
		wantTarget *aggregations.Target
		wantForm   contracts.TargetForm
	}

	testCases := []testCase{{
		"returns target if form is valid",
		contracts.TargetForm{
			Name:      contracts.Field{},
			Questions: contracts.Fields{},
		},
		&aggregations.Target{},
		contracts.TargetForm{
			Name:      contracts.Field{},
			Questions: contracts.Fields{},
		},
	}, {
		"returns nil if form is invalid and validates the form",
		contracts.TargetForm{},
		nil,
		contracts.TargetForm{
			Name: contracts.Field{
				Errors: contracts.Errors{validators.ErrEmpty},
			},
			Questions: contracts.Fields{contracts.Field{
				Errors: contracts.Errors{validators.ErrEmpty},
			}},
		},
	}}

	asm := builders.TargetAssembler{}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			gotTarget := asm.Target(&tc.form)
			assert.Equal(t, tc.wantTarget, gotTarget)
			assert.EqualExportedValues(t, tc.wantForm, tc.form)
		})
	}
}

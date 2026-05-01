package builders_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/builders"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

func TestTargetAssembler(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc       string
		form       contracts.TargetForm
		ownerID    valueobjects.ID
		wantTarget *aggregations.Target
		wantForm   contracts.TargetForm
	}

	ownerID := valueobjects.ID(69)
	testCases := []testCase{{
		"returns target if form is valid",
		contracts.TargetForm{
			Name: contracts.Field{Value: "name"},
			Questions: contracts.Fields{{
				Value: "kek1",
			}, {
				Value: "kek2",
			}},
		},
		ownerID,
		&aggregations.Target{
			Name:  "name",
			Owner: entities.User{ID: ownerID},
			Questions: valueobjects.Questions{{
				Text: "kek1",
			}, {
				Text: "kek2",
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: "name"},
			Questions: contracts.Fields{{
				Value: "kek1",
			}, {
				Value: "kek2",
			}},
		},
	}, {
		"returns nil if form is invalid and validates the form",
		contracts.TargetForm{},
		ownerID,
		nil,
		contracts.TargetForm{
			Name: contracts.Field{
				Errors: contracts.Errors{validators.ErrEmpty},
			},
			Questions: contracts.Fields{contracts.Field{
				Errors: contracts.Errors{validators.ErrEmpty},
			}},
		},
	}, {
		"returns nil if questions are wrong",
		contracts.TargetForm{
			Name: contracts.Field{Value: "name"},
			Questions: contracts.Fields{{
				Value: "kek1",
			}, {
				Value: "kek1",
			}},
		},
		ownerID,
		nil,
		contracts.TargetForm{
			Name: contracts.Field{Value: "name"},
			Questions: contracts.Fields{{
				Value: "kek1",
			}, {
				Value: "kek1",
				Errors: contracts.Errors{
					validators.ErrRepeatedQuestion,
				},
			}},
		},
	}}

	builder := builders.TargetBuilder{
		validators.TargetFormValidator{TargetRepo: targetRepo},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			gotTarget, err := builder.Target(ctx, &tc.form, tc.ownerID)
			require.NoError(t, err)
			assert.Equal(t, tc.wantTarget, gotTarget)
			assert.EqualExportedValues(t, tc.wantForm, tc.form)
		})
	}
}

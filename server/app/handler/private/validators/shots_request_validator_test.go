package validators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/validators"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

func TestShotsValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		shots       contracts.ShotsRequest
		result      validators.ValidationResult
	}

	testCases := []testCase{{
		"should be valid if everything is 0",
		contracts.ShotsRequest{
			Shots: []valueobjects.Shot{{X: 0, Y: 0}},
		},
		validators.ValidationResult{Errors: nil},
	}, {
		"should be valid if everything is on range (0; 100)",
		contracts.ShotsRequest{
			Shots: []valueobjects.Shot{{X: 69, Y: 69}},
		},
		validators.ValidationResult{Errors: nil},
	}, {
		"should be valid if everything is 100",
		contracts.ShotsRequest{
			Shots: []valueobjects.Shot{{X: 100, Y: 100}},
		},
		validators.ValidationResult{Errors: nil},
	}, {
		"should not be valid if coordinate is less than 0",
		contracts.ShotsRequest{
			Shots: []valueobjects.Shot{{X: -1, Y: 0}},
		},
		validators.ValidationResult{
			Errors: []error{validators.ErrShotBadCoordinate},
		},
	}, {
		"should not be valid if both coordinates are less than 0",
		contracts.ShotsRequest{
			Shots: []valueobjects.Shot{{X: -1, Y: -1}},
		},
		validators.ValidationResult{
			Errors: []error{
				validators.ErrShotBadCoordinate,
				validators.ErrShotBadCoordinate,
			}},
	}, {
		"should not be valid if x > 100 and y < 0",
		contracts.ShotsRequest{
			Shots: []valueobjects.Shot{{X: 101, Y: -1}},
		},
		validators.ValidationResult{
			Errors: []error{
				validators.ErrShotBadCoordinate,
				validators.ErrShotBadCoordinate,
			}},
	}, {
		"should not be valid if x < 0 and y > 100",
		contracts.ShotsRequest{
			Shots: []valueobjects.Shot{{X: -1, Y: 101}},
		},
		validators.ValidationResult{
			Errors: []error{
				validators.ErrShotBadCoordinate,
				validators.ErrShotBadCoordinate,
			}},
	}}

	v := validators.ShotsRequestValidator{}
	for _, tc := range testCases {
		res := v.Validate(tc.shots)
		assert.Len(t, res.Errors, len(tc.result.Errors))

		for i := range res.Errors {
			got := res.Errors[i]
			want := tc.result.Errors[i]
			assert.ErrorIs(t, got, want)
		}
	}
}

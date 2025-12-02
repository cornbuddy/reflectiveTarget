package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/domain/errors"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

func TestShotsValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		shots       ShotsRequest
		result      ValidationResult
	}

	testCases := []testCase{{
		"should be valid if everything is 0",
		ShotsRequest{
			Shots: []valueobjects.Shot{{X: 0, Y: 0}},
		},
		ValidationResult{Errors: nil},
	}, {
		"should be valid if everything is on range (0; 100)",
		ShotsRequest{
			Shots: []valueobjects.Shot{{X: 69, Y: 69}},
		},
		ValidationResult{Errors: nil},
	}, {
		"should be valid if everything is 100",
		ShotsRequest{
			Shots: []valueobjects.Shot{{X: 100, Y: 100}},
		},
		ValidationResult{Errors: nil},
	}, {
		"should not be valid if coordinate is less than 0",
		ShotsRequest{
			Shots: []valueobjects.Shot{{X: -1, Y: 0}},
		},
		ValidationResult{
			Errors: []error{errors.ErrShotBadCoordinate},
		},
	}, {
		"should not be valid if both coordinates are less than 0",
		ShotsRequest{
			Shots: []valueobjects.Shot{{X: -1, Y: -1}},
		},
		ValidationResult{
			Errors: []error{
				errors.ErrShotBadCoordinate,
				errors.ErrShotBadCoordinate,
			}},
	}, {
		"should not be valid if x > 100 and y < 0",
		ShotsRequest{
			Shots: []valueobjects.Shot{{X: 101, Y: -1}},
		},
		ValidationResult{
			Errors: []error{
				errors.ErrShotBadCoordinate,
				errors.ErrShotBadCoordinate,
			}},
	}, {
		"should not be valid if x < 0 and y > 100",
		ShotsRequest{
			Shots: []valueobjects.Shot{{X: -1, Y: 101}},
		},
		ValidationResult{
			Errors: []error{
				errors.ErrShotBadCoordinate,
				errors.ErrShotBadCoordinate,
			}},
	}}

	v := ShotsRequestValidator{}
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

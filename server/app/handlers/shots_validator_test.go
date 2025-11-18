package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/infra/validator"
)

func TestShotsValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		request     ShotsRequest
		result      validator.ValidationResult
	}

	testCases := []testCase{{
		"should be valid if everything is 0",
		ShotsRequest{Shots: []entities.Shot{{X: 0, Y: 0}}},
		validator.ValidationResult{Errors: nil},
	}, {
		"should be valid if everything is on range (0; 100)",
		ShotsRequest{Shots: []entities.Shot{{X: 69, Y: 69}}},
		validator.ValidationResult{Errors: nil},
	}, {
		"should be valid if everything is 100",
		ShotsRequest{Shots: []entities.Shot{{X: 100, Y: 100}}},
		validator.ValidationResult{Errors: nil},
	}, {
		"should not be valid if coordinate is less than 0",
		ShotsRequest{Shots: []entities.Shot{{X: -1, Y: 0}}},
		validator.ValidationResult{
			Errors: []error{ErrShotBadCoordinate},
		},
	}, {
		"should not be valid if both coordinates are less than 0",
		ShotsRequest{Shots: []entities.Shot{{X: -1, Y: -1}}},
		validator.ValidationResult{
			Errors: []error{
				ErrShotBadCoordinate, ErrShotBadCoordinate,
			}},
	}, {
		"should not be valid if x > 100 and y < 0",
		ShotsRequest{Shots: []entities.Shot{{X: 101, Y: -1}}},
		validator.ValidationResult{
			Errors: []error{
				ErrShotBadCoordinate, ErrShotBadCoordinate,
			}},
	}, {
		"should not be valid if x < 0 and y > 100",
		ShotsRequest{Shots: []entities.Shot{{X: -1, Y: 101}}},
		validator.ValidationResult{
			Errors: []error{
				ErrShotBadCoordinate, ErrShotBadCoordinate,
			}},
	}}

	v := ShotsRequestValidator{}
	for _, tc := range testCases {
		res := v.Validate(tc.request)
		assert.EqualValues(t, tc.result, res, tc.description)
	}
}

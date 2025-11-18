package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
)

func TestShotsValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		request     entities.ShotsRequest
		result      ValidationResult
	}

	testCases := []testCase{{
		"should be valid if everything is 0",
		entities.ShotsRequest{Shots: []entities.Shot{{X: 0, Y: 0}}},
		ValidationResult{nil},
	}, {
		"should be valid if everything is on range (0; 100)",
		entities.ShotsRequest{Shots: []entities.Shot{{X: 69, Y: 69}}},
		ValidationResult{nil},
	}, {
		"should be valid if everything is 100",
		entities.ShotsRequest{Shots: []entities.Shot{{X: 100, Y: 100}}},
		ValidationResult{nil},
	}, {
		"should not be valid if coordinate is less than 0",
		entities.ShotsRequest{Shots: []entities.Shot{{X: -1, Y: 0}}},
		ValidationResult{[]error{ErrShotBadCoordinate}},
	}, {
		"should not be valid if both coordinates are less than 0",
		entities.ShotsRequest{Shots: []entities.Shot{{X: -1, Y: -1}}},
		ValidationResult{
			[]error{
				ErrShotBadCoordinate, ErrShotBadCoordinate,
			}},
	}, {
		"should not be valid if x > 100 and y < 0",
		entities.ShotsRequest{Shots: []entities.Shot{{X: 101, Y: -1}}},
		ValidationResult{
			[]error{
				ErrShotBadCoordinate, ErrShotBadCoordinate,
			}},
	}, {
		"should not be valid if x < 0 and y > 100",
		entities.ShotsRequest{Shots: []entities.Shot{{X: -1, Y: 101}}},
		ValidationResult{
			[]error{
				ErrShotBadCoordinate, ErrShotBadCoordinate,
			}},
	}}

	validator := ShotsRequestValidator{}
	for _, tc := range testCases {
		res := validator.Validate(tc.request)
		assert.EqualValues(t, tc.result, res, tc.description)
	}
}

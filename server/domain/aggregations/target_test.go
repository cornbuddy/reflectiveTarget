package aggregations_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	aggr "github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
)

func TestTargetStringer(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc   string
		target aggr.Target
		want   string
	}

	testCases := []testCase{{
		"defaults",
		aggr.Target{},
		"{ID: 0, Name: '', " +
			"Owner: {ID: 0, Username: ''}, " +
			"Questions: [], " +
			"Shots (count): 0}",
	}, {
		"name and id",
		aggr.Target{ID: 69, Name: "kek"},
		"{ID: 69, Name: 'kek', " +
			"Owner: {ID: 0, Username: ''}, " +
			"Questions: [], " +
			"Shots (count): 0}",
	}, {
		"shots",
		aggr.Target{Shots: vo.Shots{{X: 1, Y: 1}, {X: 2, Y: 2}}},
		"{ID: 0, Name: '', " +
			"Owner: {ID: 0, Username: ''}, " +
			"Questions: [], " +
			"Shots (count): 2}",
	}, {
		"owner",
		aggr.Target{Owner: entities.User{ID: 69, Username: "kek"}},
		"{ID: 0, Name: '', " +
			"Owner: {ID: 69, Username: 'kek'}, " +
			"Questions: [], " +
			"Shots (count): 0}",
	}, {
		"questions",
		aggr.Target{
			Questions: vo.Questions{{
				ID: 1, Text: "1",
			}, {
				ID: 2, Text: "2",
			}},
		},
		"{ID: 0, Name: '', " +
			"Owner: {ID: 0, Username: ''}, " +
			"Questions: [{ID: 1, Text: '1'}, {ID: 2, Text: '2'}], " +
			"Shots (count): 0}",
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			got := tc.target.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

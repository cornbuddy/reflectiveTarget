package sessiondata_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
)

var (
	ctx = context.TODO()
)

func TestMake(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc string
		ctx  context.Context
		want SessionData
	}

	testData := SessionData{true, 69, "kek"}
	testCases := []testCase{{
		"should return zero value by default",
		ctx,
		SessionData{},
	}, {
		"should return explicitly set logger",
		context.WithValue(ctx, SessionDataCtx, &testData),
		testData,
	}}

	for _, tc := range testCases {
		got := Read(tc.ctx)
		assert.EqualExportedValues(t, tc.want, got, tc.desc)
	}
}

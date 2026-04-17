package contracts_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

func TestFieldsMethods(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc   string
		fields contracts.Fields
		valid  bool
	}

	err := errors.New("kek")
	testCases := []testCase{{
		"empty fields are valid",
		contracts.Fields{},
		true,
	}, {
		"all valid fields are valid",
		contracts.Fields{{
			Value: "",
		}, {
			Value: "",
		}},
		true,
	}, {
		"single field with error is invalid",
		contracts.Fields{{
			Errors: contracts.Errors{err},
		}},
		false,
	}, {
		"errored field in the end makes fields invalid",
		contracts.Fields{{
			Value: "kek",
		}, {
			Errors: contracts.Errors{err},
		}},
		false,
	}, {
		"errored field in the beginning makes fields invalid",
		contracts.Fields{{
			Errors: contracts.Errors{err},
		}, {
			Value: "kek",
		}},
		false,
	}, {
		"errored field in the middle makes fields invalid",
		contracts.Fields{{
			Value: "kek",
		}, {
			Errors: contracts.Errors{err},
		}, {
			Value: "kek",
		}},
		false,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.valid, tc.fields.AreValid())
			assert.Equal(t, !tc.valid, tc.fields.AreInvalid())
		})
	}
}

func TestFieldMethods(t *testing.T) {
	t.Parallel()

	field := contracts.Field{}
	assert.Empty(t, field.Errors)
	assert.True(t, field.IsValid())
	assert.False(t, field.IsInvalid())

	err := errors.New("kek")
	field.AddError(err)
	assert.Contains(t, field.Errors, err)
	assert.Len(t, field.Errors, 1)
	assert.False(t, field.IsValid())
	assert.True(t, field.IsInvalid())
}

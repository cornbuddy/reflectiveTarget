package contracts_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
)

func TestFieldsStringer(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc   string
		fields contracts.Fields
		want   string
	}

	testCases := []testCase{{
		"empty slice",
		contracts.Fields{},
		"[]",
	}, {
		"default field",
		contracts.Fields{contracts.Field{}},
		"[{ID: 0, Value: '', Errors: []}]",
	}, {
		"multiple non-default fields",
		contracts.Fields{
			contracts.Field{ID: 0, Value: kek1},
			contracts.Field{ID: 69, Value: kek2},
		},
		"[{ID: 0, Value: 'kek1', Errors: []}, {ID: 69, Value: 'kek2', Errors: []}]",
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			got := tc.fields.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestFieldStringer(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc  string
		field contracts.Field
		want  string
	}

	testCases := []testCase{{
		"default struct",
		contracts.Field{},
		"{ID: 0, Value: '', Errors: []}",
	}, {
		"errors",
		contracts.Field{
			Errors: contracts.Errors{err1, err2},
		},
		"{ID: 0, Value: '', Errors: ['kek1', 'kek2']}",
	}, {
		"ID",
		contracts.Field{ID: 69},
		"{ID: 69, Value: '', Errors: []}",
	}, {
		"value",
		contracts.Field{Value: kek},
		"{ID: 0, Value: 'kek', Errors: []}",
	}, {
		"all together",
		contracts.Field{
			ID:     420,
			Value:  "kekeke",
			Errors: contracts.Errors{err1, err2},
		},
		"{ID: 420, Value: 'kekeke', Errors: ['kek1', 'kek2']}",
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			got := tc.field.String()
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestErrors(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc string
		errs contracts.Errors
		want string
	}

	testCases := []testCase{{
		"single item",
		contracts.Errors{err},
		kek,
	}, {
		"multiple items",
		contracts.Errors{err1, err2},
		"kek1\nkek2",
	}, {
		"empty list",
		contracts.Errors{},
		"",
	}, {
		"nil",
		nil,
		"",
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, tc.errs.Error())
		})
	}
}

func TestFieldsMethods(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc   string
		fields contracts.Fields
		valid  bool
	}

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
			Value: kek,
		}, {
			Errors: contracts.Errors{err},
		}},
		false,
	}, {
		"errored field in the beginning makes fields invalid",
		contracts.Fields{{
			Errors: contracts.Errors{err},
		}, {
			Value: kek,
		}},
		false,
	}, {
		"errored field in the middle makes fields invalid",
		contracts.Fields{{
			Value: kek,
		}, {
			Errors: contracts.Errors{err},
		}, {
			Value: kek,
		}},
		false,
	}}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

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

	field.AddError(err)
	assert.Contains(t, field.Errors, err)
	assert.Len(t, field.Errors, 1)
	assert.False(t, field.IsValid())
	assert.True(t, field.IsInvalid())
}

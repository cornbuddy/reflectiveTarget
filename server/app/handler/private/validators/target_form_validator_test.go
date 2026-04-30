package validators

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestTargetFormValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		desc     string
		form     contracts.TargetForm
		wantForm contracts.TargetForm
		wantRes  bool
	}

	longTargetName := utils.MakeRandomString(MaxTargetNameLen + 1)
	longQuestion := utils.MakeRandomString(MaxQuestionLen + 1)
	allowedAmountOfQuestions := makeRandomFields(MaxNumOfQuestions)
	notSoLongTargetName := utils.MakeRandomString(MaxTargetNameLen)
	notSoLongQuestion := utils.MakeRandomString(MaxQuestionLen)
	aLotOfQuestions := makeRandomFields(MaxNumOfQuestions + 1)
	aLotOfQuestionsWithError := slices.Clone(aLotOfQuestions)
	aLotOfQuestionsWithError[len(aLotOfQuestionsWithError)-1].
		AddError(ErrExcessiveQuestion)

	testCases := []testCase{{
		"should reject if fields are empty",
		contracts.TargetForm{},
		contracts.TargetForm{
			Name: contracts.Field{
				Errors: contracts.Errors{ErrEmpty},
			},
			Questions: []contracts.Field{{
				Errors: contracts.Errors{ErrEmpty},
			}},
		},
		false,
	}, {
		"should reject if question is empty",
		contracts.TargetForm{
			Name: contracts.Field{Value: notSoLongTargetName},
			Questions: []contracts.Field{{
				Value: "",
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: notSoLongTargetName},
			Questions: []contracts.Field{{
				Value:  "",
				Errors: contracts.Errors{ErrEmpty},
			}},
		},
		false,
	}, {
		"should reject if too much questions",
		contracts.TargetForm{
			Name:      contracts.Field{Value: notSoLongTargetName},
			Questions: aLotOfQuestions,
		},
		contracts.TargetForm{
			Name:      contracts.Field{Value: notSoLongTargetName},
			Questions: aLotOfQuestionsWithError,
		},
		false,
	}, {
		"should reject if target name is too long",
		contracts.TargetForm{
			Name:      contracts.Field{Value: longTargetName},
			Questions: allowedAmountOfQuestions,
		},
		contracts.TargetForm{
			Name: contracts.Field{
				Value:  longTargetName,
				Errors: contracts.Errors{ErrTooLongTargetName},
			},
			Questions: allowedAmountOfQuestions,
		},
		false,
	}, {

		"should reject if question name is too long",
		contracts.TargetForm{
			Name: contracts.Field{Value: notSoLongTargetName},
			Questions: contracts.Fields{{
				Value: longQuestion,
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: notSoLongTargetName},
			Questions: contracts.Fields{{
				Value:  longQuestion,
				Errors: contracts.Errors{ErrTooLongQuestion},
			}},
		},
		false,
	}, {
		"should reject if question text is too long",
		contracts.TargetForm{
			Name: contracts.Field{Value: notSoLongTargetName},
			Questions: contracts.Fields{{
				Value: longQuestion,
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: notSoLongTargetName},
			Questions: contracts.Fields{{
				Value:  longQuestion,
				Errors: contracts.Errors{ErrTooLongQuestion},
			}},
		},
		false,
	}, {
		"should reject if questions are repeated",
		contracts.TargetForm{
			Name: contracts.Field{Value: notSoLongTargetName},
			Questions: contracts.Fields{{
				Value: notSoLongQuestion,
			}, {
				Value: notSoLongQuestion,
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: notSoLongTargetName},
			Questions: contracts.Fields{{
				Value: notSoLongQuestion,
			}, {
				Value:  notSoLongQuestion,
				Errors: contracts.Errors{ErrRepeatedQuestion},
			}},
		},
		false,
	}, {
		"should be valid with max number of questions",
		contracts.TargetForm{
			Name:      contracts.Field{Value: notSoLongTargetName},
			Questions: allowedAmountOfQuestions,
		},
		contracts.TargetForm{
			Name:      contracts.Field{Value: notSoLongTargetName},
			Questions: allowedAmountOfQuestions,
		},
		true,
	}, {
		"should be valid if lengths are maxed",
		contracts.TargetForm{
			Name: contracts.Field{Value: notSoLongTargetName},
			Questions: contracts.Fields{{
				Value: notSoLongQuestion,
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: notSoLongTargetName},
			Questions: contracts.Fields{{
				Value: notSoLongQuestion,
			}},
		},
		true,
	}}

	v := TargetFormValidator{}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := v.Validate(&tc.form)
			assert.Equal(t, tc.wantRes, got)
			assert.EqualExportedValues(t, tc.wantForm, tc.form)
		})
	}
}

func makeRandomFields(amount int) contracts.Fields {
	fields := make(contracts.Fields, amount)
	for i := range amount {
		fields[i] = contracts.Field{
			Value: utils.MakeRandomString(5),
		}
	}

	return fields
}

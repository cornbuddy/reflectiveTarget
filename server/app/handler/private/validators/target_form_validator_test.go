package validators

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

func TestTargetFormValidator(t *testing.T) {
	t.Parallel()

	target, owner, err := insertTestTarget(db)
	require.NoError(t, err)

	type testCase struct {
		desc     string
		ownerID  valueobjects.ID
		targetID valueobjects.ID
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

	newTargetID := valueobjects.ID(0)
	testCases := []testCase{{
		"should reject if fields are empty",
		owner.ID,
		newTargetID,
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
		owner.ID,
		newTargetID,
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
		owner.ID,
		newTargetID,
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
		owner.ID,
		newTargetID,
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
		owner.ID,
		newTargetID,
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
		owner.ID,
		newTargetID,
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
		owner.ID,
		newTargetID,
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
		"should reject if user already has target with this name",
		owner.ID,
		newTargetID,
		contracts.TargetForm{
			Name: contracts.Field{Value: target.Name},
			Questions: contracts.Fields{{
				Value: notSoLongQuestion,
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{
				Value: target.Name,
				Errors: contracts.Errors{
					ErrTargetAlreadyExists,
				},
			},
			Questions: contracts.Fields{{
				Value: notSoLongQuestion,
			}},
		},
		false,
	}, {
		"should be valid if another user has target with the same name",
		69,
		newTargetID,
		contracts.TargetForm{
			Name: contracts.Field{Value: target.Name},
			Questions: contracts.Fields{{
				Value: notSoLongQuestion,
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: target.Name},
			Questions: contracts.Fields{{
				Value: notSoLongQuestion,
			}},
		},
		true,
	}, {
		"should be valid if changing questions for existing target",
		owner.ID,
		target.ID,
		contracts.TargetForm{
			Name: contracts.Field{Value: target.Name},
			Questions: contracts.Fields{{
				Value: notSoLongQuestion,
			}},
		},
		contracts.TargetForm{
			Name: contracts.Field{Value: target.Name},
			Questions: contracts.Fields{{
				Value: notSoLongQuestion,
			}},
		},
		true,
	}, {
		"should be valid with max number of questions",
		owner.ID,
		newTargetID,
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
		owner.ID,
		newTargetID,
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

	v := TargetFormValidator{repositories.TargetRepo{DB: db}}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := v.Validate(ctx, &tc.form, tc.ownerID, tc.targetID)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, got)
			assert.EqualExportedValues(t, tc.wantForm, tc.form)
		})
	}
}

func insertTestTarget(db *sql.DB) (
	*aggregations.Target, *entities.User, error,
) {

	user, err := entities.NewUser(utils.MakeRandomString(5), "kek")
	if err != nil {
		return nil, nil, err
	}

	if err := utils.InsertUser(db, user); err != nil {
		return nil, nil, err
	}

	target := &aggregations.Target{
		Name:  utils.MakeRandomString(5),
		Owner: *user,
		Questions: valueobjects.Questions{{
			Text: utils.MakeRandomString(5),
		}, {
			Text: utils.MakeRandomString(5),
		}},
	}
	if err := utils.InsertTarget(db, target); err != nil {
		return nil, nil, err
	}

	return target, user, nil
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

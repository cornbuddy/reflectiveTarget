package repositories_test

import (
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	"github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type TargetCommandsTest struct {
	repo  *repositories.TargetRepo
	owner *entities.User
}

func (s *TargetCommandsTest) SaveShouldUpdateFieldsWhenChanged(t *testgroup.T) {
	want := makeTarget(*s.owner)
	t.Require.NoError(s.repo.Save(ctx, &want))

	wantName, wantQuestion := "changed", "changed as well"
	want.Name = wantName
	want.Questions[0].Text = wantQuestion
	t.Require.NoError(s.repo.Save(ctx, &want))
	t.Equal(wantName, want.Name)
	t.Equal(wantQuestion, want.Questions[0].Text)

	got, err := s.repo.Get(ctx, want.ID)
	t.Require.NoError(err)
	t.EqualValues(want, *got)
}

func (s *TargetCommandsTest) SaveShouldBeIdempotent(t *testgroup.T) {
	first := makeTarget(*s.owner)
	t.Require.NoError(s.repo.Save(ctx, &first))

	second, err := testutils.DeepCopy(first)
	t.Require.NoError(err)

	t.Require.NoError(s.repo.Save(ctx, second))
	t.EqualExportedValues(first, *second)
}

func (s *TargetCommandsTest) PreGroup(t *testgroup.T) {
	owner, err := entities.NewUser("user1", "password")
	t.Require.NoError(err)
	t.Require.NoError(testutils.InsertUser(db, owner))

	s.repo = &repositories.TargetRepo{db}
	s.owner = owner
}

func TestTargetQueries(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetQueriesTest))
}

func makeTarget(owner entities.User) aggregations.Target {
	return aggregations.Target{
		Name:  utils.MakeRandomString(5),
		Owner: owner,
		Questions: vo.Questions{{
			Text: utils.MakeRandomString(10),
		}, {
			Text: utils.MakeRandomString(10),
		}},
		Shots: valueobjects.Shots{},
	}
}

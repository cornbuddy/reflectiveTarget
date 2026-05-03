package repositories_test

import (
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/domain/entities"
	vo "github.com/cornbuddy/reflectiveTarget/server/domain/valueobjects"
	"github.com/cornbuddy/reflectiveTarget/server/infra/repositories"
	testutils "github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type TargetReadsTest struct {
	targets aggregations.Targets
	repo    *repositories.TargetRepo
	// owns 1 targets
	owner1 *entities.User
	// owns 2 targets
	owner2 *entities.User
}

func (s *TargetReadsTest) GetShouldReturnNilIfTargetNotExist(t *testgroup.T) {
	got, err := s.repo.Get(ctx, 69)
	t.Require.NoError(err)
	t.Nil(got)
}

func (s *TargetReadsTest) GetShouldReturnTargetIfExist(t *testgroup.T) {
	want := s.targets[0]
	got, err := s.repo.Get(ctx, want.ID)
	t.Require.NoError(err)
	t.NotNil(got)
	t.EqualValues(want, *got)
}

func (s *TargetReadsTest) ListShouldReturnEmptyListIfNoUser(t *testgroup.T) {
	targets, err := s.repo.ListTargetsOfUser(ctx, 69)
	t.Require.NoError(err)
	t.Empty(targets)
}

func (s *TargetReadsTest) ListShouldReturnTargetsIfExist(t *testgroup.T) {
	got1, err := s.repo.ListTargetsOfUser(ctx, s.owner1.ID)
	t.Require.NoError(err)
	t.Len(got1, 1)

	got2, err := s.repo.ListTargetsOfUser(ctx, s.owner2.ID)
	t.Require.NoError(err)
	t.Len(got2, 2)

	for i, got := range append(got1, got2...) {
		want := s.targets[i]
		t.Equal(want.Name, got.Name)
		t.Equal(want.ID, got.ID)
		t.Empty(got.Questions)
		t.Empty(got.Shots)
		t.Empty(got.Owner)
	}
}

func (s *TargetReadsTest) PreGroup(t *testgroup.T) {
	owner1, err := entities.NewUser("user1", "password")
	t.Require.NoError(err)

	owner2, err := entities.NewUser("user2", "password")
	t.Require.NoError(err)

	users := entities.Users{*owner1, *owner2}
	t.Require.NoError(testutils.InsertUsers(db, users))

	questions := vo.Questions{
		{Text: "kek1?"},
		{Text: "kek2?"},
	}
	shots := vo.Shots{
		{X: 1, Y: 100},
		{X: 100, Y: 1},
	}
	targets := aggregations.Targets{aggregations.Target{
		Name:      "test1",
		Owner:     users[0],
		Questions: append(vo.Questions{}, questions...),
		Shots:     append(vo.Shots{}, shots...),
	}, {
		Name:      "test2",
		Owner:     users[1],
		Questions: append(vo.Questions{}, questions...),
		Shots:     append(vo.Shots{}, shots...),
	}, {
		Name:      "test3",
		Owner:     users[1],
		Questions: append(vo.Questions{}, questions...),
		Shots:     append(vo.Shots{}, shots...),
	}}
	t.Require.NoError(testutils.InsertTargets(db, targets))

	s.owner1 = &users[0]
	s.owner2 = &users[1]
	s.repo = &repositories.TargetRepo{db}
	s.targets = targets
}

func TestTargetRepo(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetReadsTest))
}

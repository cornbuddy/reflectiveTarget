package handlers

import (
	"io"
	"net/http"
	"testing"

	"github.com/bloomberg/go-testgroup"

	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

type TargetsSuite struct {
	session        []*http.Cookie
	ownedTargets   aggregations.Targets
	foreignTargets aggregations.Targets
}

func (s *TargetsSuite) ShouldListTargetsForOwner(t *testgroup.T) {
	r := utils.MakeRequestWithCookies(
		"", http.MethodGet, "/targets", router, nil, s.session...,
	)
	t.Equal(http.StatusOK, r.StatusCode)

	data, err := io.ReadAll(r.Body)
	t.Require.NoError(err)

	t.Cleanup(func() {
		t.Require.NoError(r.Body.Close())
	})

	body := string(data)
	for _, ot := range s.ownedTargets {
		t.Contains(body, ot.ID)
		t.Contains(body, ot.Name)
	}

	for _, ft := range s.foreignTargets {
		t.NotContains(body, ft.ID)
		t.NotContains(body, ft.Name)
	}
}

func (*TargetsSuite) PreGroup(t *testgroup.T) {}

func TestTargetsHandler(t *testing.T) {
	t.Parallel()

	testgroup.RunInParallel(t, new(TargetsSuite))
}

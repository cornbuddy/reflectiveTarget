package middlewares

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/bloomberg/go-testgroup"
	"github.com/gorilla/mux"

	"github.com/cornbuddy/reflectiveTarget/server/infra/daos"
	"github.com/cornbuddy/reflectiveTarget/server/test/utils"
)

var (
	ctx       = context.TODO()
	emptyStub = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {},
	)

	store daos.SessionStore
	mw    Middleware
)

type IsAuthenticatedSuite struct {
	authenticatedToken string
	handler            http.HandlerFunc
}

func (s *IsAuthenticatedSuite) PreGroup(t *testgroup.T) {
	token := "kekeke"
	t.Require.NoError(store.SaveSession(ctx, token, true))

	s.authenticatedToken = token
	s.handler = chain(
		emptyStub,
		mw.IsAuthenticated,
		mw.PutSessionDataToContext,
	).ServeHTTP
}

func TestMain(m *testing.M) {
	cleanup, cache, err := utils.SetupCache(ctx)
	if err != nil {
		log.Fatalf("failed to setup cache: %v", err)
	}

	defer func() {
		if err := cleanup(); err != nil {
			log.Fatalf("failed to cleanup cache: %v", err)
		}
	}()

	store = daos.SessionStore{Cache: cache}
	mw = Middleware{SessionStore: store}

	os.Exit(m.Run())
}

func chain(mux http.Handler, mwf ...mux.MiddlewareFunc) http.Handler {
	for _, mw := range mwf {
		mux = mw(mux)
	}

	return mux
}

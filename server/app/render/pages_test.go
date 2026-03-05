package render_test

import (
	"context"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cornbuddy/reflectiveTarget/server/app/constants"
	"github.com/cornbuddy/reflectiveTarget/server/app/render"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
)

func TestPageShouldContainText(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		render      render.RenderFunc
		ctx         context.Context
		data        any
		contains    []string
	}

	emptyCtx := context.TODO()
	authorizedCtx := context.WithValue(
		context.TODO(),
		constants.AuthenticatedCtx,
		toPtr(true),
	)

	testCases := []testCase{{
		"non authorized client should see correct navbar",
		render.Layout.Index,
		emptyCtx,
		nil,
		[]string{"Login", "Signup"},
	}, {
		"authorized client should see correct navbar",
		render.Layout.Index,
		authorizedCtx,
		nil,
		[]string{"Targets", "Logout"},
	}, {
		"signup page should contain form",
		render.View.Signup,
		emptyCtx,
		nil,
		[]string{"Signup", "<form hx-post=\"/signup\"", "</form>"},
	}, {
		"login page should contain form",
		render.View.Login,
		emptyCtx,
		nil,
		[]string{"Login", "<form hx-post=\"/login\"", "</form>"},
	}, {
		"targets page should allow to create new target",
		render.View.Targets,
		authorizedCtx,
		nil,
		[]string{"Targets", "New target", "href=\"/targets/new\""},
	}, {
		"targets page should have links to targets",
		render.View.Targets,
		authorizedCtx,
		map[string]any{
			"Targets": aggregations.Targets{{
				ID:   1,
				Name: "kek1?",
			}, {
				ID:   69,
				Name: "kek69?",
			}},
		},
		[]string{
			"href=\"/targets/1\">kek1?<",
			"href=\"/targets/69\">kek69?<",
		},
	}}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		tc.render(tc.ctx, w, tc.data)
		raw, err := io.ReadAll(w.Body)
		require.NoError(t, err, tc.description)

		body := string(raw)
		for _, token := range tc.contains {
			assert.Contains(t, body, token, tc.description)
		}
	}
}

var (
	anyType = reflect.TypeFor[struct{}]()
)

func TestViewRendererShouldNotContainFullPage(t *testing.T) {
	t.Parallel()

	value := reflect.ValueOf(&render.View)
	for i := 0; i < value.NumMethod(); i++ {
		method := value.Method(i)
		require.True(t, method.IsValid())

		w := httptest.NewRecorder()
		args := []reflect.Value{
			reflect.ValueOf(context.TODO()),
			reflect.ValueOf(w),
			reflect.New(anyType).Elem(),
		}
		method.Call(args)
		raw, err := io.ReadAll(w.Body)
		require.NoError(t, err)

		body := string(raw)
		assert.NotContains(t, body, "<!DOCTYPE html>")
		assert.NotContains(t, body, "</html>")
	}
}

func TestLayoutRendererShouldContainFullPage(t *testing.T) {
	t.Parallel()

	value := reflect.ValueOf(&render.Layout)
	for i := 0; i < value.NumMethod(); i++ {
		method := value.Method(i)
		require.True(t, method.IsValid())

		w := httptest.NewRecorder()
		args := []reflect.Value{
			reflect.ValueOf(context.TODO()),
			reflect.ValueOf(w),
			reflect.New(anyType).Elem(),
		}
		method.Call(args)
		raw, err := io.ReadAll(w.Body)
		require.NoError(t, err)

		body := string(raw)
		assert.Contains(t, body, "<!DOCTYPE html>")
		assert.Contains(t, body, "</html>")
	}
}

func toPtr[A any](obj A) *A {
	return &obj
}

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
)

func TestAuthzHandlerShouldRenderFormsOnGet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		description string
		render      render.RenderFunc
		ctx         context.Context
		w           *httptest.ResponseRecorder
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
		httptest.NewRecorder(),
		nil,
		[]string{"Login", "Signup"},
	}, {
		"authorized client should see correct navbar",
		render.Layout.Index,
		authorizedCtx,
		httptest.NewRecorder(),
		nil,
		[]string{"Logout"},
	}, {
		"signup page should contain form",
		render.View.Signup,
		emptyCtx,
		httptest.NewRecorder(),
		nil,
		[]string{
			"Signup", "<form hx-post=\"/signup\"", "</form>",
		},
	}, {
		"login page should contain form",
		render.View.Login,
		emptyCtx,
		httptest.NewRecorder(),
		nil,
		[]string{
			"Login", "<form hx-post=\"/login\"", "</form>",
		},
	}}

	for _, tc := range testCases {
		tc.render(tc.ctx, tc.w, tc.data)
		raw, err := io.ReadAll(tc.w.Body)
		require.NoError(t, err, tc.description)

		body := string(raw)
		for _, token := range tc.contains {
			assert.Contains(t, body, token, tc.description)
		}
	}
}

func TestViewRendererShouldNotContainFullPage(t *testing.T) {
	t.Parallel()

	type stub struct{}
	anyType := reflect.TypeOf(stub{})

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

	type stub struct{}
	anyType := reflect.TypeOf(stub{})

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

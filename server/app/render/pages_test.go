package render_test

import (
	"context"
	"io"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/cornbuddy/reflectiveTarget/server/app/render"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestAuthzHandlerShouldRenderFormsOnGet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		render   render.RenderFunc
		contains []string
	}

	testCases := []testCase{{
		render.View.Signup,
		[]string{"Signup", "</form>"},
	}, {
		render.View.Login,
		[]string{"Login", "</form>"},
	}}

	for _, tc := range testCases {
		ctx := context.TODO()
		w := httptest.NewRecorder()
		tc.render(ctx, w, nil)
		raw, err := io.ReadAll(w.Body)
		require.NoError(t, err)

		body := string(raw)

		for _, token := range tc.contains {
			assert.Contains(t, body, token)
		}
	}
}

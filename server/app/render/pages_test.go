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

// func TestAuthzHandlerShouldRenderFormsOnGet(t *testing.T) {
// 	t.Parallel()
//
// 	type testCase struct {
// 		url      string
// 		contains string
// 	}
//
// 	testCases := []testCase{{
// 		url:      "/signup",
// 		contains: "Signup",
// 	}, {
// 		url:      "/login",
// 		contains: "Login",
// 	}}
//
// 	for _, tc := range testCases {
// 		get := http.MethodGet
// 		res := utils.MakeRequest("", get, tc.url, router, nil)
// 		ct := "text/html; charset=utf-8"
// 		assert.Equal(t, ct, res.Header.Get("Content-Type"))
// 		assert.Equal(t, http.StatusOK, res.StatusCode)
//
// 		data, err := io.ReadAll(res.Body)
// 		assert.NoError(t, err)
//
// 		t.Cleanup(func() { res.Body.Close() })
//
// 		body := string(data)
// 		assert.Contains(t, body, tc.contains)
// 		assert.Contains(t, body, "</form>")
// 	}
// }

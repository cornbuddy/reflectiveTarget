package render

import (
	"context"
	"net/http"

	"github.com/abiosoft/mold"
)

type RenderFunc func(context.Context, http.ResponseWriter, any)

type engine struct {
	Engine mold.Engine
}

type contextData struct {
	IsAuthorized bool
}

type viewData struct {
	contextData
	any
}

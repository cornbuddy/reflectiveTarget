package render

import (
	"context"
	"net/http"

	"github.com/abiosoft/mold"
)

type RenderFunc func(context.Context, http.ResponseWriter, ViewData)

type ViewData map[string]any

type engine struct {
	Engine mold.Engine
}

package render

import (
	"context"
	"net/http"
)

func (e engine) Index(ctx context.Context, w http.ResponseWriter, data any) {
	e.Engine.Render(w, "views/index.tmpl", data)
}

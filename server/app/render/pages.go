package render

import (
	"context"
	"net/http"
)

func (e engine) Signup(ctx context.Context, w http.ResponseWriter, data any) {
	e.Engine.Render(w, "views/signup.tmpl", data)
}

func (e engine) Login(ctx context.Context, w http.ResponseWriter, data any) {
	e.Engine.Render(w, "views/login.tmpl", data)
}

func (e engine) Index(ctx context.Context, w http.ResponseWriter, data any) {
	e.Engine.Render(w, "views/index.tmpl", data)
}

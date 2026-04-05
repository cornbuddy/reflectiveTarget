package render

import (
	"context"
	"embed"
	"net/http"

	"github.com/abiosoft/mold"

	"github.com/cornbuddy/reflectiveTarget/server/app/forms"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
)

type Render interface {
	Index(context.Context, http.ResponseWriter)
	Login(context.Context, http.ResponseWriter, LoginData)
	Signup(context.Context, http.ResponseWriter, SignupData)
	Targets(context.Context, http.ResponseWriter, TargetsData)
}

// renders view with layout (eg navbar, head, header, footer, etc)
var Layout = makeRender("layout.tmpl")

// renders view only
var View = makeRender("empty-layout.tmpl")

type TargetsData struct {
	aggregations.Targets
}

type SignupData struct {
	forms.SignupForm
}

type LoginData struct {
	forms.LoginForm
}

type viewData map[string]any

type engine struct {
	mold.Engine
}

//go:embed templates
var templates embed.FS

func makeRender(layout string) Render {
	return engine{mold.Must(mold.New(templates, mold.With(
		mold.WithRoot("templates"),
		mold.WithLayout(layout),
	)))}
}

package render

import (
	"context"
	"embed"
	"net/http"

	"github.com/abiosoft/mold"
	"github.com/go-task/slim-sprig/v3"

	"github.com/cornbuddy/reflectiveTarget/server/app/handler/private/contracts"
	"github.com/cornbuddy/reflectiveTarget/server/app/sessiondata"
	"github.com/cornbuddy/reflectiveTarget/server/domain/aggregations"
)

type Render interface {
	Index(context.Context, http.ResponseWriter)
	Login(context.Context, http.ResponseWriter, LoginData)
	Signup(context.Context, http.ResponseWriter, SignupData)
	Targets(context.Context, http.ResponseWriter, TargetsData)
	UpdateTarget(context.Context, http.ResponseWriter, TargetData)
}

// renders view with layout (eg navbar, head, header, footer, etc)
var Layout = makeRender("layout.tmpl")

// renders view only
var View = makeRender("empty-layout.tmpl")

type TargetData struct {
	contracts.TargetForm
}

type TargetsData struct {
	aggregations.Targets
}

type SignupData struct {
	contracts.SignupForm
}

type LoginData struct {
	contracts.LoginForm
}

type viewData struct {
	SessionData sessiondata.SessionData
	TargetsData TargetsData
	TargetData  TargetData
	SignupData  SignupData
	LoginData   LoginData
}

type engine struct {
	mold.Engine
}

//go:embed templates
var templates embed.FS

func makeRender(layout string) Render {
	return engine{mold.Must(mold.New(templates, mold.With(
		mold.WithRoot("templates"),
		mold.WithFuncMap(sprig.FuncMap()),
		mold.WithLayout(layout),
	)))}
}

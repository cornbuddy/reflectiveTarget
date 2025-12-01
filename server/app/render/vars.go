package render

import (
	"context"
	"embed"
	"net/http"

	"github.com/abiosoft/mold"
)

type RenderFunc func(context.Context, http.ResponseWriter, any)

var Layout = engine{layoutEngine}

var View = engine{viewsEngine}

type engine struct {
	Engine mold.Engine
}

//go:embed templates
var dir embed.FS

var layoutEngine = mold.Must(mold.New(dir, mold.With(
	mold.WithRoot("templates"),
	mold.WithLayout("layout.tmpl"),
)))

var viewsEngine = mold.Must(mold.New(dir, mold.With(
	mold.WithRoot("templates"),
)))

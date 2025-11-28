package render

import (
	"embed"

	"github.com/abiosoft/mold"
)

var Layout = engine{layoutEngine}

var View = engine{viewsEngine}

type engine struct {
	mold.Engine
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

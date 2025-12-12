package render

import (
	"embed"

	"github.com/abiosoft/mold"
)

var Layout = engine{layoutEngine}

var View = engine{viewsEngine}

//go:embed templates
var dir embed.FS

var layoutEngine = mold.Must(mold.New(dir, mold.With(
	mold.WithRoot("templates"),
	mold.WithLayout("layout.tmpl"),
)))

var viewsEngine = mold.Must(mold.New(dir, mold.With(
	mold.WithRoot("templates"),
	mold.WithLayout("empty-layout.tmpl"),
)))

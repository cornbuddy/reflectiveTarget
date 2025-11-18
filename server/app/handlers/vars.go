package handlers

import (
	"embed"

	"github.com/abiosoft/mold"
)

//go:embed templates
var dir embed.FS
var layout = mold.Must(mold.New(dir, mold.With(
	mold.WithRoot("templates"),
	mold.WithLayout("layout.tmpl"),
)))

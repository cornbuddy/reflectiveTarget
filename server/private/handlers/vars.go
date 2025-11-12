package handlers

import (
	"embed"

	"github.com/abiosoft/mold"
)

//go:embed templates
var dir embed.FS
var options = mold.With(
	mold.WithRoot("templates"),
	mold.WithLayout("layout.tmpl"),
)
var engine = mold.Must(mold.New(dir, options))

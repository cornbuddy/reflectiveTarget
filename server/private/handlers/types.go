package handlers

import (
	"database/sql"
	"embed"

	"github.com/abiosoft/mold"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
	"github.com/cornbuddy/reflectiveTarget/server/private/validators"
)

//go:embed templates
var dir embed.FS
var options = mold.With(
	mold.WithRoot("templates"),
	mold.WithLayout("layout.tmpl"),
)
var engine = mold.Must(mold.New(dir, options))

type AuthzHandler struct {
	daos.UserDao
}

type HealthHandler struct {
	*sql.DB
}

type IndexHandler struct{}

type ShotsHandler struct {
	daos.ShotsDao
	Validator validators.ShotsRequestValidator
}

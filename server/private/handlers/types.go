package handlers

import (
	"database/sql"
	"embed"

	"github.com/abiosoft/mold"

	"github.com/cornbuddy/reflectiveTarget/server/private/daos"
)

//go:embed templates
var dir embed.FS
var options = mold.With(
	mold.WithRoot("templates"),
	mold.WithLayout("layout.tmpl"),
)
var engine = mold.Must(mold.New(dir, options))

type AuthzRouter struct {
	daos.UserDao
}

type HealthRouter struct {
	*sql.DB
}

type IndexRouter struct{}

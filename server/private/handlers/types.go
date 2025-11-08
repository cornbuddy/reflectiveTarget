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

type HealthResponse struct {
	Connected   bool `json:"connected"`
	Connections int  `json:"connections"`
}

type HealthRouter struct {
	*sql.DB
}

type IndexRouter struct{}

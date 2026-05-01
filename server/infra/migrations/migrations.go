package migrations

import _ "embed"

//go:embed tables.sql
var InitQuery string

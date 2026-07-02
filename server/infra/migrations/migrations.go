package migrations

import _ "embed"

//go:embed init-db.sql
var InitDBQuery string

//go:embed init-session-store.sql
var InitSessionStoreQuery string

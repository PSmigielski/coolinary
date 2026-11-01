// Package migrations embeds the SQL migrations (goose format).
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS

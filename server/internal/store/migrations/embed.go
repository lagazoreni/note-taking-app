package migrations

import "embed"

// Files contains append-only SQL migrations shipped with the server binary.
//
//go:embed *.sql
var Files embed.FS

package sqlite

import (
	"context"
	"database/sql"

	"noted.local/noted/internal/importexport"
)

// ApplyArchive is the storage boundary used by the import service. The archive
// package validates before this transaction is started.
func ApplyArchive(ctx context.Context, db *sql.DB, archive importexport.Archive, resolutions map[string]string) (map[string]int, error) {
	return importexport.Apply(ctx, db, archive, resolutions)
}

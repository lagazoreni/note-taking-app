package testsupport

import (
	"context"
	"database/sql"
	"testing"

	"noted.local/noted/internal/store/sqlite"
)

func OpenMigratedDB(t *testing.T) (*sql.DB, func()) {
	db, cleanup := OpenTempDB(t)
	if err := sqlite.RunMigrations(context.Background(), db); err != nil {
		cleanup()
		t.Fatal(err)
	}
	return db, cleanup
}

func AssertFTSContains(db *sql.DB, table, term string) error {
	var count int
	return db.QueryRow("SELECT count(*) FROM "+table+" WHERE "+table+" MATCH ?", term).Scan(&count)
}

package sqlite_test

import (
	"context"
	"testing"

	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestMigrationsCreateSchemaAndAreIdempotent(t *testing.T) {
	ctx := context.Background()
	db, cleanup := testsupport.OpenTempDB(t)
	defer cleanup()
	if err := sqlite.RunMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if err := sqlite.RunMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"workspaces", "notes", "questions", "note_questions", "schema_migrations"} {
		var count int
		if err := db.QueryRowContext(ctx, "SELECT count(*) FROM sqlite_master WHERE type IN ('table','view') AND name = ?", table).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Fatalf("table %s missing", table)
		}
	}
}

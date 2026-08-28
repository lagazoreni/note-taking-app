package sqlite_test

import (
	"context"
	"testing"

	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestDiagnosticsAndFTSRebuild(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	ctx := context.Background()
	if err := sqlite.IntegrityCheck(ctx, db); err != nil {
		t.Fatal(err)
	}
	if err := sqlite.ForeignKeyCheck(ctx, db); err != nil {
		t.Fatal(err)
	}
	if err := sqlite.RebuildFTS(ctx, db); err != nil {
		t.Fatal(err)
	}
}

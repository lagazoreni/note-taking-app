package sqlite_test

import (
	"context"
	"testing"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/importexport"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestExportRoundTripIncludesCanonicalData(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(platformTime()), nil)
	ws, err := store.CreateWorkspace(context.Background(), "Work")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.CreateQuestion(context.Background(), domain.QuestionWrite{WorkspaceID: ws.ID, QuestionText: "Portable?"}); err != nil {
		t.Fatal(err)
	}
	archive, err := importexport.Export(context.Background(), db, "test", platform.NewFixedClock(platformTime()))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := importexport.ValidateArchive(archive, 10<<20); err != nil {
		t.Fatal(err)
	}
}

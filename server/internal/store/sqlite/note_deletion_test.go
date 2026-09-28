package sqlite_test

import (
	"context"
	"testing"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestAtomicNoteDeletionCanKeepQuestionUnlinked(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(platformTime()), nil)
	ws, _ := store.CreateWorkspace(context.Background(), "Work")
	q, _ := store.CreateQuestion(context.Background(), domain.QuestionWrite{WorkspaceID: ws.ID, QuestionText: "Keep"})
	n, _ := store.CreateNote(context.Background(), domain.NoteWrite{WorkspaceID: ws.ID, Title: "Delete", BodyMarkdown: "{{question:" + q.ID + "}}", QuestionLinks: []domain.NoteQuestionWrite{{QuestionID: q.ID, Position: 0, DisplayMode: domain.DisplayCollapsed}}})
	if err := store.DeleteNote(context.Background(), n.ID, n.Version, map[string]string{q.ID: "keep_unlinked"}, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := store.GetQuestion(context.Background(), q.ID); err != nil {
		t.Fatalf("question lost: %v", err)
	}
}

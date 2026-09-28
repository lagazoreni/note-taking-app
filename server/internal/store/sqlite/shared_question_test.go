package sqlite_test

import (
	"context"
	"testing"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestSharedQuestionHasOneCanonicalRecordAndVersionConflicts(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(platformTime()), platform.NewIDs("11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222", "33333333-3333-4333-8333-333333333333", "44444444-4444-4444-8444-444444444444"))
	ws, _ := store.CreateWorkspace(context.Background(), "Work")
	q, _ := store.CreateQuestion(context.Background(), domain.QuestionWrite{WorkspaceID: ws.ID, QuestionText: "Original?"})
	first, _ := store.CreateNote(context.Background(), domain.NoteWrite{WorkspaceID: ws.ID, Title: "First", BodyMarkdown: "one"})
	second, _ := store.CreateNote(context.Background(), domain.NoteWrite{WorkspaceID: ws.ID, Title: "Second", BodyMarkdown: "two"})
	if err := store.LinkExistingQuestion(context.Background(), first.ID, q.ID, domain.DisplayCollapsed, 0); err != nil {
		t.Fatal(err)
	}
	if err := store.LinkExistingQuestion(context.Background(), second.ID, q.ID, domain.DisplayExpanded, 0); err != nil {
		t.Fatal(err)
	}
	value, _ := store.GetQuestion(context.Background(), q.ID)
	if len(value.LinkedNotes) != 2 {
		t.Fatalf("linked notes=%d", len(value.LinkedNotes))
	}
	write := domain.QuestionWrite{WorkspaceID: ws.ID, QuestionText: "Changed?", Status: value.Status, Priority: value.Priority}
	if _, err := store.UpdateQuestion(context.Background(), q.ID, write, value.Version-1); err == nil {
		t.Fatal("stale update accepted")
	}
}

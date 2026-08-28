package service

import (
	"context"
	"testing"
	"time"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestSharedQuestionLinkingRejectsDuplicateAndCrossWorkspace(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	ids := platform.NewIDs("11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222", "33333333-3333-4333-8333-333333333333", "44444444-4444-4444-8444-444444444444")
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)), ids)
	one, _ := store.CreateWorkspace(context.Background(), "One")
	two, _ := store.CreateWorkspace(context.Background(), "Two")
	q, _ := store.CreateQuestion(context.Background(), domain.QuestionWrite{WorkspaceID: one.ID, QuestionText: "Shared?"})
	n, _ := store.CreateNote(context.Background(), domain.NoteWrite{WorkspaceID: one.ID, Title: "A", BodyMarkdown: "source"})
	svc := NewSharedQuestionService(store)
	if err := svc.Link(context.Background(), n.ID, q.ID, domain.DisplayCollapsed, 0); err != nil {
		t.Fatal(err)
	}
	if err := svc.Link(context.Background(), n.ID, q.ID, domain.DisplayCollapsed, 0); err == nil {
		t.Fatal("duplicate link accepted")
	}
	other, _ := store.CreateNote(context.Background(), domain.NoteWrite{WorkspaceID: two.ID, Title: "B", BodyMarkdown: "source"})
	if err := svc.Link(context.Background(), other.ID, q.ID, domain.DisplayCollapsed, 0); err == nil {
		t.Fatal("cross-workspace link accepted")
	}
}

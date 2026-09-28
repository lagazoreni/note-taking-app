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

func TestDeletionPreviewRequiresExplicitSinglyLinkedDecisions(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)), nil)
	ws, _ := store.CreateWorkspace(context.Background(), "Work")
	q, _ := store.CreateQuestion(context.Background(), domain.QuestionWrite{WorkspaceID: ws.ID, QuestionText: "Keep?"})
	n, _ := store.CreateNote(context.Background(), domain.NoteWrite{WorkspaceID: ws.ID, Title: "Delete me", BodyMarkdown: "{{question:" + q.ID + "}}", QuestionLinks: []domain.NoteQuestionWrite{{QuestionID: q.ID, Position: 0, DisplayMode: domain.DisplayCollapsed}}})
	svc := NewNoteDeletionService(store, platform.NewFixedClock(time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)))
	preview, err := svc.Preview(context.Background(), n.ID, n.Version)
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.SinglyLinkedQuestions) != 1 {
		t.Fatal("missing singly linked question")
	}
	if err := svc.Execute(context.Background(), n.ID, preview.Token, n.Version, nil, nil); err == nil {
		t.Fatal("delete without a decision accepted")
	}
}

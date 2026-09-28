package sqlite_test

import (
	"context"
	"testing"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestLifecycleUpdateAndFilters(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(platformTime()), nil)
	ws, _ := store.CreateWorkspace(context.Background(), "Work")
	q, _ := store.CreateQuestion(context.Background(), domain.QuestionWrite{WorkspaceID: ws.ID, QuestionText: "Answer me", Status: domain.StatusUnanswered, Priority: domain.PriorityNone})
	answer := "An answer"
	updated, err := store.UpdateQuestion(context.Background(), q.ID, domain.QuestionWrite{WorkspaceID: ws.ID, QuestionText: q.QuestionText, AnswerMarkdown: &answer, Status: domain.StatusAnswered, Priority: q.Priority}, q.Version)
	if err != nil || updated.Status != domain.StatusAnswered {
		t.Fatalf("answer update: %+v %v", updated, err)
	}
	if updated.AnswerMarkdown == nil {
		t.Fatal("answer was lost")
	}
	active, err := store.ListQuestions(context.Background(), ws.ID, []domain.QuestionStatus{domain.StatusUnanswered, domain.StatusInProgress, domain.StatusDeferred})
	if err != nil {
		t.Fatal(err)
	}
	if len(active) != 0 {
		t.Fatalf("active count=%d", len(active))
	}
}

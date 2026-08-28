package sqlite_test

import (
	"context"
	"testing"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestQuestionQueryTopicMatchDoesNotDuplicate(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(platformTime()), nil)
	ws, _ := store.CreateWorkspace(context.Background(), "Work")
	topic, _ := store.CreateTopic(context.Background(), ws.ID, "Research")
	q, _ := store.CreateQuestion(context.Background(), domain.QuestionWrite{WorkspaceID: ws.ID, QuestionText: "One question"})
	n1, _ := store.CreateNote(context.Background(), domain.NoteWrite{WorkspaceID: ws.ID, Title: "A", BodyMarkdown: "{{question:" + q.ID + "}}", TopicID: &topic.ID, QuestionLinks: []domain.NoteQuestionWrite{{QuestionID: q.ID, Position: 0, DisplayMode: domain.DisplayCollapsed}}})
	_, _ = store.CreateNote(context.Background(), domain.NoteWrite{WorkspaceID: ws.ID, Title: "B", BodyMarkdown: "{{question:" + q.ID + "}}", QuestionLinks: []domain.NoteQuestionWrite{{QuestionID: q.ID, Position: 0, DisplayMode: domain.DisplayCollapsed}}})
	_ = n1
	page, err := store.QueryQuestions(context.Background(), domain.QuestionQuery{WorkspaceID: ws.ID, TopicID: &topic.ID, PageSize: 50})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 {
		t.Fatalf("expected one canonical result, got %d", len(page.Items))
	}
}

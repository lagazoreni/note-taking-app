package sqlite_test

import (
	"context"
	"testing"
	"time"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestCaptureForeignKeysAndFTS(t *testing.T) {
	db, cleanup := testsupport.OpenTempDB(t)
	defer cleanup()
	if err := sqlite.RunMigrations(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(platformTime()), platform.NewIDs(
		"11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222", "33333333-3333-4333-8333-333333333333"))
	workspace, err := store.CreateWorkspace(context.Background(), "Work")
	if err != nil {
		t.Fatal(err)
	}
	question, err := store.CreateQuestion(context.Background(), domain.QuestionWrite{WorkspaceID: workspace.ID, QuestionText: "Why?", Status: domain.StatusUnanswered, Priority: domain.PriorityNone})
	if err != nil {
		t.Fatal(err)
	}
	note, err := store.CreateNote(context.Background(), domain.NoteWrite{WorkspaceID: workspace.ID, Title: "Research", BodyMarkdown: "A question {{question:22222222-2222-4222-8222-222222222222}}", QuestionLinks: []domain.NoteQuestionWrite{{QuestionID: question.ID, Position: 0, DisplayMode: domain.DisplayCollapsed}}})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.ReplaceNoteQuestions(context.Background(), note.ID, []domain.NoteQuestionWrite{{QuestionID: question.ID, Position: 0, DisplayMode: domain.DisplayCollapsed}}); err != nil {
		t.Fatal(err)
	}
	var count int
	if err := db.QueryRow("SELECT count(*) FROM note_questions").Scan(&count); err != nil || count != 1 {
		t.Fatalf("link count=%d err=%v", count, err)
	}
	if err := testsupport.AssertFTSContains(db, "questions_fts", "why"); err != nil {
		t.Fatal(err)
	}
}

func TestQuestionKindAndActiveListExcludesAnnotations(t *testing.T) {
	db, cleanup := testsupport.OpenTempDB(t)
	defer cleanup()
	if err := sqlite.RunMigrations(context.Background(), db); err != nil {
		t.Fatal(err)
	}
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(platformTime()), platform.NewIDs(
		"aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa",
		"bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb",
		"cccccccc-cccc-4ccc-8ccc-cccccccccccc",
	))
	workspace, err := store.CreateWorkspace(context.Background(), "Work")
	if err != nil {
		t.Fatal(err)
	}
	question, err := store.CreateQuestion(context.Background(), domain.QuestionWrite{
		WorkspaceID: workspace.ID, QuestionText: "Why?", Status: domain.StatusUnanswered, Priority: domain.PriorityNone,
	})
	if err != nil {
		t.Fatal(err)
	}
	annotation, err := store.CreateQuestion(context.Background(), domain.QuestionWrite{
		WorkspaceID: workspace.ID, QuestionText: "Margin note", Status: domain.StatusUnanswered, Priority: domain.PriorityNone, Kind: domain.KindAnnotation,
	})
	if err != nil {
		t.Fatal(err)
	}
	if question.Kind != domain.KindQuestion {
		t.Fatalf("default kind=%q", question.Kind)
	}
	if annotation.Kind != domain.KindAnnotation {
		t.Fatalf("annotation kind=%q", annotation.Kind)
	}
	page, err := store.QueryQuestions(context.Background(), domain.QuestionQuery{
		WorkspaceID: workspace.ID, Kind: domain.KindQuestion,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 || page.Items[0].ID != question.ID {
		t.Fatalf("active list should exclude annotations: %+v", page.Items)
	}
	updated, err := store.UpdateLifecycle(context.Background(), annotation.ID, domain.StatusInProgress, nil, annotation.Version)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Kind != domain.KindAnnotation {
		t.Fatalf("lifecycle dropped kind: %q", updated.Kind)
	}
}

func platformTime() (t time.Time) { return time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC) }

package sqlite_test

import (
	"context"
	"testing"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestFTSSearchAndScopedResults(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(platformTime()), nil)
	ws, _ := store.CreateWorkspace(context.Background(), "Work")
	_, _ = store.CreateQuestion(context.Background(), domain.QuestionWrite{WorkspaceID: ws.ID, QuestionText: "Find the lighthouse"})
	results, err := store.Search(context.Background(), "light-house", domain.SearchEverything, domain.SearchAll, nil, 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) == 0 {
		t.Fatal("FTS search did not find question")
	}
}

func TestFTSSearchFiltersNotesByTagWithoutDuplicates(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	store := sqlite.NewCaptureStore(db, platform.NewFixedClock(platformTime()), nil)
	ws, err := store.CreateWorkspace(context.Background(), "Work")
	if err != nil {
		t.Fatal(err)
	}
	firstTag, err := store.CreateTag(context.Background(), ws.ID, "first")
	if err != nil {
		t.Fatal(err)
	}
	secondTag, err := store.CreateTag(context.Background(), ws.ID, "second")
	if err != nil {
		t.Fatal(err)
	}
	for _, tags := range [][]string{{firstTag.ID, secondTag.ID}, {secondTag.ID}} {
		if _, err := store.CreateNote(context.Background(), domain.NoteWrite{
			WorkspaceID:   ws.ID,
			Title:         "Lighthouse note",
			BodyMarkdown:  "The lighthouse is useful.",
			QuestionLinks: []domain.NoteQuestionWrite{},
			TagIDs:        tags,
		}); err != nil {
			t.Fatal(err)
		}
	}
	results, err := store.SearchWithTag(context.Background(), "lighthouse", domain.SearchNotes, domain.SearchCurrent, &ws.ID, &secondTag.ID, 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Fatalf("expected two tagged notes, got %d", len(results))
	}
	seen := map[string]bool{}
	for _, result := range results {
		if seen[result.ID] {
			t.Fatalf("duplicate note result %s", result.ID)
		}
		seen[result.ID] = true
	}
	firstResults, err := store.Search(context.Background(), "lighthouse", domain.SearchNotes, domain.SearchCurrent, &ws.ID, &firstTag.ID, 50)
	if err != nil {
		t.Fatal(err)
	}
	if len(firstResults) != 1 {
		t.Fatalf("expected one first-tagged note, got %d", len(firstResults))
	}
}

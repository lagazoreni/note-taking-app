package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestSearchAndQuestionQueryValidation(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	h := NewAPI(db)
	req := httptest.NewRequest("GET", "/api/v1/search?q=x&contentScope=everything&workspaceScope=current", nil)
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	if resp.Code != 422 {
		t.Fatalf("expected validation status, got %d", resp.Code)
	}
}

func TestSearchAcceptsTagFilterForNotes(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	store := sqlite.NewCaptureStore(db, platform.SystemClock{}, nil)
	workspace, err := store.CreateWorkspace(context.Background(), "Work")
	if err != nil {
		t.Fatal(err)
	}
	tag, err := store.CreateTag(context.Background(), workspace.ID, "project")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.CreateNote(context.Background(), domain.NoteWrite{
		WorkspaceID:   workspace.ID,
		Title:         "Tagged note",
		BodyMarkdown:  "A lighthouse project.",
		QuestionLinks: []domain.NoteQuestionWrite{},
		TagIDs:        []string{tag.ID},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := store.CreateNote(context.Background(), domain.NoteWrite{
		WorkspaceID:   workspace.ID,
		Title:         "Other note",
		BodyMarkdown:  "A lighthouse project.",
		QuestionLinks: []domain.NoteQuestionWrite{},
		TagIDs:        []string{},
	}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/search?q=lighthouse&contentScope=notes&workspaceScope=current&workspaceId="+workspace.ID+"&tagId="+tag.ID, nil)
	resp := httptest.NewRecorder()
	NewAPI(db).ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected search success, got %d", resp.Code)
	}
	var page struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&page); err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 {
		t.Fatalf("expected one tagged result, got %d", len(page.Items))
	}
}

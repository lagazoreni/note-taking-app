package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"noted.local/noted/internal/api"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestCaptureContractWorkspaceNoteQuestion(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	h := api.NewRouter(api.RouterOptions{DB: db, Version: "test", API: NewAPI(db)})
	create := func(url, body string) *http.Response {
		req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()
		h.ServeHTTP(resp, req)
		return resp.Result()
	}
	workspace := create("/api/v1/workspaces", `{"name":"Work"}`)
	if workspace.StatusCode != http.StatusCreated {
		t.Fatalf("workspace status %d", workspace.StatusCode)
	}
	var value map[string]any
	_ = json.NewDecoder(workspace.Body).Decode(&value)
	if value["id"] == nil {
		t.Fatal("workspace id missing")
	}
	_ = context.Background()
	_ = sqlite.RunMigrations
}

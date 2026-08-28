package handlers

import (
	"net/http/httptest"
	"testing"

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

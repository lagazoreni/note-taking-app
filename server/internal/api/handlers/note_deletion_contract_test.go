package handlers

import (
	"net/http/httptest"
	"testing"

	"noted.local/noted/internal/testsupport"
)

func TestDeletionPreviewRequiresVersion(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	req := httptest.NewRequest("POST", "/api/v1/notes/x/deletion-preview", nil)
	out := httptest.NewRecorder()
	NewAPI(db).ServeHTTP(out, req)
	if out.Code != 422 {
		t.Fatalf("status=%d", out.Code)
	}
}

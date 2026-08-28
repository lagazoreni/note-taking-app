package handlers

import (
	"net/http/httptest"
	"testing"

	"noted.local/noted/internal/testsupport"
)

func TestExportEndpointReturnsArchiveContentType(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	req := httptest.NewRequest("POST", "/api/v1/export", nil)
	out := httptest.NewRecorder()
	NewAPI(db).ServeHTTP(out, req)
	if out.Code != 200 {
		t.Fatalf("status=%d", out.Code)
	}
	if out.Header().Get("Content-Type") != "application/zip" {
		t.Fatalf("content type=%q", out.Header().Get("Content-Type"))
	}
}

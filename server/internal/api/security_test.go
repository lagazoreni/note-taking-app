package api_test

import (
	"context"
	"testing"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestSecurityInputsAreRejectedOrRedacted(t *testing.T) {
	if _, err := domain.ParseQuestionDirectives("{{question:<script>alert(1)</script>}}"); err == nil {
		t.Fatal("unsafe directive accepted")
	}
	if err := domain.ValidateNoteText("title", string(make([]byte, 5_000_001))); err == nil {
		t.Fatal("oversized note accepted")
	}
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	if err := sqlite.CheckFTS5(context.Background(), db); err != nil {
		t.Fatal(err)
	}
}

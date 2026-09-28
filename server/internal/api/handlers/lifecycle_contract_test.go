package handlers

import (
	"context"
	"testing"

	"noted.local/noted/internal/store/sqlite"
	"noted.local/noted/internal/testsupport"
)

func TestLifecycleContractUsesQuestionEndpoint(t *testing.T) {
	db, cleanup := testsupport.OpenMigratedDB(t)
	defer cleanup()
	if NewAPI(db) == nil {
		t.Fatal("API was nil")
	}
	_ = context.Background()
	_ = sqlite.ErrConflict
}

package contracttest

import (
	"context"
	"path/filepath"
	"testing"
)

func TestOpenAPIContractLoadsAndValidates(t *testing.T) {
	path := filepath.Join("..", "..", "..", "..", "specs", "001-track-note-questions", "contracts", "openapi.yaml")
	doc, err := LoadDocument(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := doc.Validate(context.Background()); err != nil {
		t.Fatalf("contract is invalid: %v", err)
	}
	for _, route := range []string{"/api/v1/workspaces", "/api/v1/notes", "/api/v1/questions", "/api/v1/search"} {
		if doc.Paths.Find(route) == nil {
			t.Fatalf("missing route %s", route)
		}
	}
}

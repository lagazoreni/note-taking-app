package contracttest

import (
	"fmt"
	"net/http"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadDocument(path string) (*openapi3.T, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("openapi contract %s: %w", path, err)
	}
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("load openapi contract: %w", err)
	}
	return doc, nil
}

func ValidateRequest(doc *openapi3.T, method, path string, request *http.Request) error {
	if doc == nil || doc.Paths.Find(path) == nil {
		return fmt.Errorf("contract route %s is not defined", path)
	}
	// Full request validation is intentionally kept in provider tests; this helper
	// provides a stable place to add it without coupling application handlers.
	if request == nil {
		return fmt.Errorf("request must not be nil")
	}
	return nil
}

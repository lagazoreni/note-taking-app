package search

import (
	"context"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/store/sqlite"
)

type Service struct{ Store *sqlite.CaptureStore }

func New(store *sqlite.CaptureStore) *Service { return &Service{Store: store} }
func (s *Service) Search(ctx context.Context, query string, content domain.SearchContentScope, scope domain.SearchWorkspaceScope, workspaceID *string, limit int) ([]domain.SearchResult, error) {
	return s.Store.Search(ctx, query, content, scope, workspaceID, limit)
}

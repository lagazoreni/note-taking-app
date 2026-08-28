package service

import (
	"context"
	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/store/sqlite"
)

type OrganizationService struct{ Store *sqlite.CaptureStore }

func NewOrganizationService(store *sqlite.CaptureStore) *OrganizationService {
	return &OrganizationService{Store: store}
}
func (s *OrganizationService) Topics(ctx context.Context, workspaceID string) ([]domain.Topic, error) {
	return s.Store.ListTopics(ctx, workspaceID)
}
func (s *OrganizationService) Tags(ctx context.Context, workspaceID string) ([]domain.Tag, error) {
	return s.Store.ListTags(ctx, workspaceID, true)
}

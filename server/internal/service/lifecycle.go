package service

import (
	"context"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/store/sqlite"
)

type LifecycleService struct{ Store *sqlite.CaptureStore }

func NewLifecycleService(store *sqlite.CaptureStore) *LifecycleService {
	return &LifecycleService{Store: store}
}
func (s *LifecycleService) Update(ctx context.Context, id string, status domain.QuestionStatus, answer *string, version int64) (domain.Question, error) {
	return s.Store.UpdateLifecycle(ctx, id, status, answer, version)
}
func (s *LifecycleService) Active(ctx context.Context, workspaceID string) ([]domain.Question, error) {
	return s.Store.ListQuestions(ctx, workspaceID, []domain.QuestionStatus{domain.StatusUnanswered, domain.StatusInProgress, domain.StatusDeferred})
}
func (s *LifecycleService) Answered(ctx context.Context, workspaceID string) ([]domain.Question, error) {
	return s.Store.ListQuestions(ctx, workspaceID, []domain.QuestionStatus{domain.StatusAnswered})
}

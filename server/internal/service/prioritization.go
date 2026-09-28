package service

import (
	"context"
	"time"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/store/sqlite"
)

type PrioritizationService struct{ Store *sqlite.CaptureStore }

func NewPrioritizationService(store *sqlite.CaptureStore) *PrioritizationService {
	return &PrioritizationService{Store: store}
}
func (s *PrioritizationService) Query(ctx context.Context, query domain.QuestionQuery) (domain.QuestionPage, error) {
	return s.Store.QueryQuestions(ctx, query)
}
func (s *PrioritizationService) EvaluateReminders(ctx context.Context, now time.Time) ([]domain.Reminder, error) {
	return s.Store.EvaluateReminders(ctx, now)
}

package service

import (
	"context"
	"fmt"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/store/sqlite"
)

type SharedQuestionService struct{ Store *sqlite.CaptureStore }

func NewSharedQuestionService(store *sqlite.CaptureStore) *SharedQuestionService {
	return &SharedQuestionService{Store: store}
}
func (s *SharedQuestionService) Link(ctx context.Context, noteID, questionID string, mode domain.DisplayMode, position int) error {
	if noteID == "" || questionID == "" {
		return fmt.Errorf("note and question are required")
	}
	return s.Store.LinkExistingQuestion(ctx, noteID, questionID, mode, position)
}
func (s *SharedQuestionService) Unlink(ctx context.Context, noteID, questionID string) error {
	return s.Store.UnlinkExistingQuestion(ctx, noteID, questionID)
}
func (s *SharedQuestionService) Reorder(ctx context.Context, noteID string, links []domain.NoteQuestionWrite) error {
	return s.Store.ReorderQuestionLinks(ctx, noteID, links)
}
func (s *SharedQuestionService) Search(ctx context.Context, workspaceID, query string) ([]domain.Question, error) {
	return s.Store.FindQuestions(ctx, workspaceID, query, 50)
}

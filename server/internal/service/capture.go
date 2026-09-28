package service

import (
	"context"
	"errors"
	"fmt"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/store/sqlite"
)

type CaptureService struct{ Store *sqlite.CaptureStore }

func NewCaptureService(store *sqlite.CaptureStore) *CaptureService {
	return &CaptureService{Store: store}
}

func (s *CaptureService) CreateWorkspace(ctx context.Context, name string) (domain.Workspace, error) {
	return s.Store.CreateWorkspace(ctx, name)
}
func (s *CaptureService) SaveNote(ctx context.Context, write domain.NoteWrite, id string, version int64) (domain.Note, error) {
	if err := write.Validate(); err != nil {
		return domain.Note{}, err
	}
	if id == "" {
		return s.Store.CreateNote(ctx, write)
	}
	return s.Store.UpdateNote(ctx, id, write, version)
}
func (s *CaptureService) SaveQuestion(ctx context.Context, write domain.QuestionWrite, id string, version int64) (domain.Question, error) {
	q := domain.Question{WorkspaceID: write.WorkspaceID, QuestionText: write.QuestionText, AnswerMarkdown: write.AnswerMarkdown, Status: write.Status, Priority: write.Priority, DueDate: write.DueDate}
	if err := q.Validate(); err != nil {
		return domain.Question{}, err
	}
	if id == "" {
		return s.Store.CreateQuestion(ctx, write)
	}
	return s.Store.UpdateQuestion(ctx, id, write, version)
}
func (s *CaptureService) ListActiveQuestions(ctx context.Context, workspaceID string) ([]domain.Question, error) {
	if workspaceID == "" {
		return nil, fmt.Errorf("workspace is required")
	}
	return s.Store.ListQuestions(ctx, workspaceID, []domain.QuestionStatus{domain.StatusUnanswered, domain.StatusInProgress, domain.StatusDeferred})
}
func IsConflict(err error) bool { return errors.Is(err, sqlite.ErrConflict) }

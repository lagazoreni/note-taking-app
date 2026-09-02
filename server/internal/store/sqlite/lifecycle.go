package sqlite

import (
	"context"

	"noted.local/noted/internal/domain"
)

func (s *CaptureStore) UpdateLifecycle(ctx context.Context, id string, status domain.QuestionStatus, answer *string, version int64) (domain.Question, error) {
	q, err := s.GetQuestion(ctx, id)
	if err != nil {
		return domain.Question{}, err
	}
	answer = domain.NormalizeAnswer(answer)
	if err := domain.ValidateTransition(q.Status, status, answer, q.DueDate); err != nil {
		return domain.Question{}, err
	}
	return s.UpdateQuestion(ctx, id, domain.QuestionWrite{WorkspaceID: q.WorkspaceID, QuestionText: q.QuestionText, Kind: q.Kind, AnswerMarkdown: answer, Status: status, Priority: q.Priority, DueDate: q.DueDate, TagIDs: q.TagIDs}, version)
}

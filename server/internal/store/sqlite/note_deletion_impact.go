package sqlite

import (
	"context"
	"noted.local/noted/internal/domain"
)

type ChildNote struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Version int64  `json:"version"`
}
type NoteDeletionImpact struct {
	Note     domain.Note
	Singly   []domain.QuestionSummary
	Multiple []domain.QuestionSummary
	Children []ChildNote
}

func (s *CaptureStore) NoteDeletionImpact(ctx context.Context, noteID string) (NoteDeletionImpact, error) {
	note, err := s.GetNote(ctx, noteID)
	if err != nil {
		return NoteDeletionImpact{}, err
	}
	rows, err := s.db.QueryContext(ctx, `SELECT nq.question_id,(SELECT count(*) FROM note_questions other WHERE other.question_id=nq.question_id AND other.note_id<>nq.note_id) FROM note_questions nq WHERE nq.note_id=?`, noteID)
	if err != nil {
		return NoteDeletionImpact{}, err
	}
	defer rows.Close()
	impact := NoteDeletionImpact{Note: note}
	for rows.Next() {
		var id string
		var others int
		if err := rows.Scan(&id, &others); err != nil {
			return impact, err
		}
		q, err := s.GetQuestion(ctx, id)
		if err != nil {
			return impact, err
		}
		summary := domain.QuestionSummary{ID: q.ID, WorkspaceID: q.WorkspaceID, QuestionText: q.QuestionText, Kind: q.Kind, Status: q.Status, Priority: q.Priority, DueDate: q.DueDate, Version: q.Version}
		if others == 0 {
			impact.Singly = append(impact.Singly, summary)
		} else {
			impact.Multiple = append(impact.Multiple, summary)
		}
	}
	if err := rows.Err(); err != nil {
		return impact, err
	}
	children, err := s.db.QueryContext(ctx, "SELECT id,title,version FROM notes WHERE parent_note_id=? ORDER BY title COLLATE NOCASE,id", noteID)
	if err != nil {
		return impact, err
	}
	defer children.Close()
	for children.Next() {
		var child ChildNote
		if err := children.Scan(&child.ID, &child.Title, &child.Version); err != nil {
			return impact, err
		}
		impact.Children = append(impact.Children, child)
	}
	return impact, children.Err()
}

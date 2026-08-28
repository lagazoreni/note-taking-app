package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"noted.local/noted/internal/domain"
)

// LinkExistingQuestion adds one relationship while keeping the question canonical.
// The caller normally updates the note body in the same save; this primitive is
// also useful to service and migration code that already owns a transaction.
func (s *CaptureStore) LinkExistingQuestion(ctx context.Context, noteID, questionID string, mode domain.DisplayMode, position int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var workspace string
	if err := tx.QueryRowContext(ctx, "SELECT workspace_id FROM notes WHERE id = ?", noteID).Scan(&workspace); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	var questionWorkspace string
	if err := tx.QueryRowContext(ctx, "SELECT workspace_id FROM questions WHERE id = ?", questionID).Scan(&questionWorkspace); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if workspace != questionWorkspace {
		return ErrWorkspaceBoundary
	}
	var count int
	if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM note_questions WHERE note_id = ? AND question_id = ?", noteID, questionID).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("question is already linked")
	}
	if err := s.insertLink(ctx, tx, noteID, workspace, domain.NoteQuestionWrite{QuestionID: questionID, DisplayMode: mode, Position: position}); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *CaptureStore) UnlinkExistingQuestion(ctx context.Context, noteID, questionID string) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM note_questions WHERE note_id = ? AND question_id = ?", noteID, questionID)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *CaptureStore) ReorderQuestionLinks(ctx context.Context, noteID string, links []domain.NoteQuestionWrite) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, "DELETE FROM note_questions WHERE note_id = ?", noteID); err != nil {
		return err
	}
	var workspace string
	if err := tx.QueryRowContext(ctx, "SELECT workspace_id FROM notes WHERE id = ?", noteID).Scan(&workspace); err != nil {
		return err
	}
	for index, link := range links {
		link.Position = index
		if err := s.insertLink(ctx, tx, noteID, workspace, link); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *CaptureStore) FindQuestions(ctx context.Context, workspaceID, query string, limit int) ([]domain.Question, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `SELECT q.id FROM questions q WHERE q.workspace_id = ? AND q.question_text LIKE ? ORDER BY q.updated_at DESC,q.id DESC LIMIT ?`, workspaceID, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Question, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		q, err := s.GetQuestion(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, q)
	}
	return result, rows.Err()
}

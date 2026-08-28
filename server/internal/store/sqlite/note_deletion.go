package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (s *CaptureStore) DeleteNote(ctx context.Context, noteID string, version int64, questionDecisions map[string]string, childDecisions map[string]string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var workspace string
	var parent sql.NullString
	var current int64
	if err := tx.QueryRowContext(ctx, "SELECT workspace_id,parent_note_id,version FROM notes WHERE id=?", noteID).Scan(&workspace, &parent, &current); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	if current != version {
		return ErrConflict
	}
	rows, err := tx.QueryContext(ctx, `SELECT nq.question_id, (SELECT count(*) FROM note_questions other WHERE other.question_id=nq.question_id AND other.note_id<>nq.note_id) FROM note_questions nq WHERE nq.note_id=?`, noteID)
	if err != nil {
		return err
	}
	type linked struct {
		id     string
		others int
	}
	var links []linked
	for rows.Next() {
		var value linked
		if err := rows.Scan(&value.id, &value.others); err != nil {
			rows.Close()
			return err
		}
		links = append(links, value)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for _, link := range links {
		if link.others == 0 {
			decision := questionDecisions[link.id]
			if decision != "keep_unlinked" && decision != "delete" {
				return fmt.Errorf("an explicit decision is required for question %s", link.id)
			}
		}
	}
	children, err := tx.QueryContext(ctx, "SELECT id FROM notes WHERE parent_note_id=?", noteID)
	if err != nil {
		return err
	}
	var childIDs []string
	for children.Next() {
		var id string
		if err := children.Scan(&id); err != nil {
			children.Close()
			return err
		}
		childIDs = append(childIDs, id)
	}
	children.Close()
	for _, childID := range childIDs {
		action := childDecisions[childID]
		if action == "" {
			return fmt.Errorf("an explicit child decision is required for note %s", childID)
		}
		if action == "delete_with_review" {
			return fmt.Errorf("nested child deletion requires its own review")
		}
		newParent := ""
		if action == "move_to_parent" && parent.Valid {
			newParent = parent.String
		}
		if action != "make_root" && action != "move_to_parent" {
			return fmt.Errorf("invalid child decision")
		}
		if _, err := tx.ExecContext(ctx, "UPDATE notes SET parent_note_id=?,updated_at=?,version=version+1 WHERE id=?", nullableString(newParent), s.now(), childID); err != nil {
			return err
		}
	}
	for _, link := range links {
		if _, err := tx.ExecContext(ctx, "DELETE FROM note_questions WHERE note_id=? AND question_id=?", noteID, link.id); err != nil {
			return err
		}
		if link.others == 0 && questionDecisions[link.id] == "delete" {
			if _, err := tx.ExecContext(ctx, "DELETE FROM question_tags WHERE question_id=?", link.id); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, "DELETE FROM questions WHERE id=?", link.id); err != nil {
				return err
			}
		}
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM note_tags WHERE note_id=?", noteID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM notes WHERE id=? AND version=?", noteID, version); err != nil {
		return err
	}
	return tx.Commit()
}
func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}

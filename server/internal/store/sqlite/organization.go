package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"noted.local/noted/internal/domain"
	"noted.local/noted/internal/platform"
)

func (s *CaptureStore) CreateTopic(ctx context.Context, workspaceID, name string) (domain.Topic, error) {
	name, err := domain.ValidateName(name, 100)
	if err != nil {
		return domain.Topic{}, err
	}
	if err := s.requireWorkspace(ctx, workspaceID); err != nil {
		return domain.Topic{}, err
	}
	id, err := s.newID()
	if err != nil {
		return domain.Topic{}, err
	}
	now := s.now()
	_, err = s.db.ExecContext(ctx, "INSERT INTO topics(id,workspace_id,name,created_at,updated_at,version) VALUES(?,?,?,?,?,1)", id, workspaceID, name, now, now)
	if err != nil {
		if isUnique(err) {
			return domain.Topic{}, ErrNameConflict
		}
		return domain.Topic{}, err
	}
	return s.GetTopic(ctx, id)
}
func (s *CaptureStore) GetTopic(ctx context.Context, id string) (domain.Topic, error) {
	var value domain.Topic
	var created, updated string
	err := s.db.QueryRowContext(ctx, "SELECT id,workspace_id,name,created_at,updated_at,version FROM topics WHERE id=?", id).Scan(&value.ID, &value.WorkspaceID, &value.Name, &created, &updated, &value.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return value, ErrNotFound
	}
	if err != nil {
		return value, err
	}
	value.CreatedAt, err = platform.ParseTimestamp(created)
	if err != nil {
		return value, err
	}
	value.UpdatedAt, err = platform.ParseTimestamp(updated)
	return value, err
}
func (s *CaptureStore) ListTopics(ctx context.Context, workspaceID string) ([]domain.Topic, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id FROM topics WHERE workspace_id=? ORDER BY name COLLATE NOCASE,id", workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Topic, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		value, err := s.GetTopic(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, rows.Err()
}
func (s *CaptureStore) UpdateTopic(ctx context.Context, id, name string, version int64) (domain.Topic, error) {
	name, err := domain.ValidateName(name, 100)
	if err != nil {
		return domain.Topic{}, err
	}
	now := s.now()
	result, err := s.db.ExecContext(ctx, "UPDATE topics SET name=?,updated_at=?,version=version+1 WHERE id=? AND version=?", name, now, id, version)
	if err != nil {
		if isUnique(err) {
			return domain.Topic{}, ErrNameConflict
		}
		return domain.Topic{}, err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		if _, e := s.GetTopic(ctx, id); errors.Is(e, ErrNotFound) {
			return domain.Topic{}, ErrNotFound
		}
		return domain.Topic{}, ErrConflict
	}
	return s.GetTopic(ctx, id)
}
func (s *CaptureStore) DeleteTopic(ctx context.Context, id string, version int64) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM topics WHERE id=? AND version=?", id, version)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		if _, e := s.GetTopic(ctx, id); errors.Is(e, ErrNotFound) {
			return ErrNotFound
		}
		return ErrConflict
	}
	return nil
}

func (s *CaptureStore) CreateTag(ctx context.Context, workspaceID, name string) (domain.Tag, error) {
	name, err := domain.ValidateName(name, 100)
	if err != nil {
		return domain.Tag{}, err
	}
	if err := s.requireWorkspace(ctx, workspaceID); err != nil {
		return domain.Tag{}, err
	}
	id, err := s.newID()
	if err != nil {
		return domain.Tag{}, err
	}
	now := s.now()
	_, err = s.db.ExecContext(ctx, "INSERT INTO tags(id,owner_workspace_id,name,created_at,updated_at,version) VALUES(?,?,?,?,?,1)", id, workspaceID, name, now, now)
	if err != nil {
		if isUnique(err) {
			return domain.Tag{}, ErrNameConflict
		}
		return domain.Tag{}, err
	}
	return s.GetTag(ctx, id)
}
func (s *CaptureStore) GetTag(ctx context.Context, id string) (domain.Tag, error) {
	var value domain.Tag
	var created, updated string
	err := s.db.QueryRowContext(ctx, "SELECT id,owner_workspace_id,name,created_at,updated_at,version FROM tags WHERE id=?", id).Scan(&value.ID, &value.OwnerWorkspaceID, &value.Name, &created, &updated, &value.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return value, ErrNotFound
	}
	if err != nil {
		return value, err
	}
	value.CreatedAt, err = platform.ParseTimestamp(created)
	if err != nil {
		return value, err
	}
	value.UpdatedAt, err = platform.ParseTimestamp(updated)
	if err != nil {
		return value, err
	}
	value.AvailableWorkspaceIDs, err = s.tagAccess(ctx, id, value.OwnerWorkspaceID)
	return value, err
}
func (s *CaptureStore) tagAccess(ctx context.Context, id, owner string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT workspace_id FROM tag_workspace_access WHERE tag_id=? ORDER BY workspace_id", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []string{owner}
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, rows.Err()
}
func (s *CaptureStore) ListTags(ctx context.Context, workspaceID string, includeShared bool) ([]domain.Tag, error) {
	query := "SELECT t.id FROM tags t WHERE t.owner_workspace_id=?"
	args := []any{workspaceID}
	if includeShared {
		query += " OR EXISTS (SELECT 1 FROM tag_workspace_access a WHERE a.tag_id=t.id AND a.workspace_id=?)"
		args = append(args, workspaceID)
	}
	query += " ORDER BY t.name COLLATE NOCASE,t.id"
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Tag, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		value, err := s.GetTag(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, rows.Err()
}
func (s *CaptureStore) UpdateTag(ctx context.Context, id, name string, version int64) (domain.Tag, error) {
	name, err := domain.ValidateName(name, 100)
	if err != nil {
		return domain.Tag{}, err
	}
	now := s.now()
	result, err := s.db.ExecContext(ctx, "UPDATE tags SET name=?,updated_at=?,version=version+1 WHERE id=? AND version=?", name, now, id, version)
	if err != nil {
		if isUnique(err) {
			return domain.Tag{}, ErrNameConflict
		}
		return domain.Tag{}, err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		if _, e := s.GetTag(ctx, id); errors.Is(e, ErrNotFound) {
			return domain.Tag{}, ErrNotFound
		}
		return domain.Tag{}, ErrConflict
	}
	return s.GetTag(ctx, id)
}
func (s *CaptureStore) SetTagAccess(ctx context.Context, id string, workspaceIDs []string, version int64, removeAssignments map[string]bool) (domain.Tag, error) {
	tag, err := s.GetTag(ctx, id)
	if err != nil {
		return domain.Tag{}, err
	}
	if tag.Version != version {
		return domain.Tag{}, ErrConflict
	}
	desired := map[string]bool{}
	for _, workspaceID := range workspaceIDs {
		if workspaceID == tag.OwnerWorkspaceID {
			continue
		}
		desired[workspaceID] = true
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Tag{}, err
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx, "SELECT workspace_id FROM tag_workspace_access WHERE tag_id=?", id)
	if err != nil {
		return domain.Tag{}, err
	}
	current := []string{}
	for rows.Next() {
		var workspaceID string
		if err := rows.Scan(&workspaceID); err != nil {
			rows.Close()
			return domain.Tag{}, err
		}
		current = append(current, workspaceID)
	}
	rows.Close()
	for _, workspaceID := range current {
		if desired[workspaceID] {
			continue
		}
		var used int
		if err := tx.QueryRowContext(ctx, "SELECT (SELECT count(*) FROM note_tags nt JOIN notes n ON n.id=nt.note_id WHERE nt.tag_id=? AND n.workspace_id=?) + (SELECT count(*) FROM question_tags qt JOIN questions q ON q.id=qt.question_id WHERE qt.tag_id=? AND q.workspace_id=?)", id, workspaceID, id, workspaceID).Scan(&used); err != nil {
			return domain.Tag{}, err
		}
		if used > 0 && !removeAssignments[workspaceID] {
			return domain.Tag{}, fmt.Errorf("tag assignments require confirmation")
		}
		if used > 0 {
			if _, err := tx.ExecContext(ctx, "DELETE FROM note_tags WHERE tag_id=? AND note_id IN (SELECT id FROM notes WHERE workspace_id=?)", id, workspaceID); err != nil {
				return domain.Tag{}, err
			}
			if _, err := tx.ExecContext(ctx, "DELETE FROM question_tags WHERE tag_id=? AND question_id IN (SELECT id FROM questions WHERE workspace_id=?)", id, workspaceID); err != nil {
				return domain.Tag{}, err
			}
		}
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM tag_workspace_access WHERE tag_id=?", id); err != nil {
		return domain.Tag{}, err
	}
	for workspaceID := range desired {
		var count int
		if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM workspaces WHERE id=?", workspaceID).Scan(&count); err != nil {
			return domain.Tag{}, err
		}
		if count != 1 {
			return domain.Tag{}, ErrWorkspaceBoundary
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO tag_workspace_access(tag_id,workspace_id,created_at) VALUES(?,?,?)", id, workspaceID, s.now()); err != nil {
			return domain.Tag{}, err
		}
	}
	if _, err := tx.ExecContext(ctx, "UPDATE tags SET version=version+1,updated_at=? WHERE id=? AND version=?", s.now(), id, version); err != nil {
		return domain.Tag{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.Tag{}, err
	}
	return s.GetTag(ctx, id)
}
func (s *CaptureStore) requireWorkspace(ctx context.Context, id string) error {
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT count(*) FROM workspaces WHERE id=?", id).Scan(&count)
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrNotFound
	}
	return nil
}

func (s *CaptureStore) AssignNoteTag(ctx context.Context, noteID, tagID string) error {
	return s.assignTag(ctx, "note_tags", "note_id", "notes", noteID, tagID)
}
func (s *CaptureStore) AssignQuestionTag(ctx context.Context, qID, tagID string) error {
	return s.assignTag(ctx, "question_tags", "question_id", "questions", qID, tagID)
}
func (s *CaptureStore) assignTag(ctx context.Context, table, entityColumn, entityTable, entityID, tagID string) error {
	var workspace string
	if err := s.db.QueryRowContext(ctx, "SELECT workspace_id FROM "+entityTable+" WHERE id=?", entityID).Scan(&workspace); err != nil {
		return err
	}
	var owner string
	if err := s.db.QueryRowContext(ctx, "SELECT owner_workspace_id FROM tags WHERE id=?", tagID).Scan(&owner); err != nil {
		return err
	}
	if owner != workspace {
		var count int
		if err := s.db.QueryRowContext(ctx, "SELECT count(*) FROM tag_workspace_access WHERE tag_id=? AND workspace_id=?", tagID, workspace).Scan(&count); err != nil {
			return err
		}
		if count != 1 {
			return ErrWorkspaceBoundary
		}
	}
	_, err := s.db.ExecContext(ctx, "INSERT OR IGNORE INTO "+table+"("+entityColumn+",tag_id,created_at) VALUES(?,?,?)", entityID, tagID, s.now())
	return err
}

func (s *CaptureStore) _organizationStringHelper(value string) string {
	return strings.TrimSpace(value)
}

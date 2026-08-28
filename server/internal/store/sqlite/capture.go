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

type IDSource interface{ Next() (string, error) }
type randomIDs struct{}

func (randomIDs) Next() (string, error) { return platform.NewID() }

type CaptureStore struct {
	db    *sql.DB
	clock platform.Clock
	ids   IDSource
}

func NewCaptureStore(db *sql.DB, clock platform.Clock, ids IDSource) *CaptureStore {
	if clock == nil {
		clock = platform.SystemClock{}
	}
	if ids == nil {
		ids = randomIDs{}
	}
	return &CaptureStore{db: db, clock: clock, ids: ids}
}
func (s *CaptureStore) now() string            { return platform.FormatTimestamp(s.clock.Now()) }
func (s *CaptureStore) newID() (string, error) { return s.ids.Next() }

func (s *CaptureStore) CreateWorkspace(ctx context.Context, name string) (domain.Workspace, error) {
	name, err := domain.ValidateName(name, 100)
	if err != nil {
		return domain.Workspace{}, err
	}
	id, err := s.newID()
	if err != nil {
		return domain.Workspace{}, err
	}
	now := s.now()
	_, err = s.db.ExecContext(ctx, "INSERT INTO workspaces(id,name,created_at,updated_at,version) VALUES(?,?,?,?,1)", id, name, now, now)
	if err != nil {
		if isUnique(err) {
			return domain.Workspace{}, ErrNameConflict
		}
		return domain.Workspace{}, fmt.Errorf("create workspace: %w", err)
	}
	return domain.Workspace{ID: id, Name: name, CreatedAt: s.clock.Now(), UpdatedAt: s.clock.Now(), Version: 1}, nil
}

func (s *CaptureStore) ListWorkspaces(ctx context.Context) ([]domain.Workspace, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id,name,created_at,updated_at,version FROM workspaces ORDER BY name COLLATE NOCASE,id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.Workspace, 0)
	for rows.Next() {
		var value domain.Workspace
		var created, updated string
		if err := rows.Scan(&value.ID, &value.Name, &created, &updated, &value.Version); err != nil {
			return nil, err
		}
		value.CreatedAt, err = platform.ParseTimestamp(created)
		if err != nil {
			return nil, err
		}
		value.UpdatedAt, err = platform.ParseTimestamp(updated)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, rows.Err()
}
func (s *CaptureStore) GetWorkspace(ctx context.Context, id string) (domain.Workspace, error) {
	var value domain.Workspace
	var created, updated string
	err := s.db.QueryRowContext(ctx, "SELECT id,name,created_at,updated_at,version FROM workspaces WHERE id=?", id).Scan(&value.ID, &value.Name, &created, &updated, &value.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Workspace{}, ErrNotFound
	}
	if err != nil {
		return domain.Workspace{}, err
	}
	value.CreatedAt, err = platform.ParseTimestamp(created)
	if err != nil {
		return value, err
	}
	value.UpdatedAt, err = platform.ParseTimestamp(updated)
	return value, err
}
func (s *CaptureStore) UpdateWorkspace(ctx context.Context, id, name string, version int64) (domain.Workspace, error) {
	name, err := domain.ValidateName(name, 100)
	if err != nil {
		return domain.Workspace{}, err
	}
	now := s.now()
	result, err := s.db.ExecContext(ctx, "UPDATE workspaces SET name=?,updated_at=?,version=version+1 WHERE id=? AND version=?", name, now, id, version)
	if err != nil {
		if isUnique(err) {
			return domain.Workspace{}, ErrNameConflict
		}
		return domain.Workspace{}, err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		if _, e := s.GetWorkspace(ctx, id); errors.Is(e, ErrNotFound) {
			return domain.Workspace{}, ErrNotFound
		}
		return domain.Workspace{}, ErrConflict
	}
	return s.GetWorkspace(ctx, id)
}

func (s *CaptureStore) CreateQuestion(ctx context.Context, write domain.QuestionWrite) (domain.Question, error) {
	q := domain.Question{WorkspaceID: write.WorkspaceID, QuestionText: write.QuestionText, Kind: write.Kind, AnswerMarkdown: write.AnswerMarkdown, Status: write.Status, Priority: write.Priority, DueDate: write.DueDate}
	if err := q.Validate(); err != nil {
		return domain.Question{}, err
	}
	id, err := s.newID()
	if err != nil {
		return domain.Question{}, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Question{}, err
	}
	defer tx.Rollback()
	var workspaceCount int
	if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM workspaces WHERE id=?", q.WorkspaceID).Scan(&workspaceCount); err != nil {
		return domain.Question{}, err
	}
	if workspaceCount != 1 {
		return domain.Question{}, ErrNotFound
	}
	now := s.now()
	if _, err := tx.ExecContext(ctx, "INSERT INTO questions(id,workspace_id,question_text,kind,answer_markdown,status,priority,due_date,created_at,updated_at,version) VALUES(?,?,?,?,?,?,?,?,?,?,1)", id, q.WorkspaceID, q.QuestionText, q.Kind, q.AnswerMarkdown, q.Status, q.Priority, q.DueDate, now, now); err != nil {
		return domain.Question{}, err
	}
	if err := s.replaceQuestionTagsTx(ctx, tx, id, q.WorkspaceID, write.TagIDs); err != nil {
		return domain.Question{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.Question{}, err
	}
	return s.GetQuestion(ctx, id)
}
func (s *CaptureStore) GetQuestion(ctx context.Context, id string) (domain.Question, error) {
	q, err := s.getQuestionBase(ctx, id)
	if err != nil {
		return domain.Question{}, err
	}
	q.LinkedNotes, err = s.linkedNotes(ctx, id)
	if err != nil {
		return q, err
	}
	q.TagIDs, err = s.questionTags(ctx, id)
	if err != nil {
		return q, err
	}
	q.Reminder, err = s.reminder(ctx, id)
	return q, err
}
func (s *CaptureStore) getQuestionBase(ctx context.Context, id string) (domain.Question, error) {
	var q domain.Question
	var answer, due sql.NullString
	var created, updated string
	err := s.db.QueryRowContext(ctx, "SELECT id,workspace_id,question_text,kind,answer_markdown,status,priority,due_date,created_at,updated_at,version FROM questions WHERE id=?", id).Scan(&q.ID, &q.WorkspaceID, &q.QuestionText, &q.Kind, &answer, &q.Status, &q.Priority, &due, &created, &updated, &q.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return q, ErrNotFound
	}
	if err != nil {
		return q, err
	}
	if answer.Valid {
		q.AnswerMarkdown = &answer.String
	}
	if due.Valid {
		q.DueDate = &due.String
	}
	q.CreatedAt, err = platform.ParseTimestamp(created)
	if err != nil {
		return q, err
	}
	q.UpdatedAt, err = platform.ParseTimestamp(updated)
	return q, err
}
func (s *CaptureStore) UpdateQuestion(ctx context.Context, id string, write domain.QuestionWrite, version int64) (domain.Question, error) {
	old, err := s.getQuestionBase(ctx, id)
	if err != nil {
		return domain.Question{}, err
	}
	q := domain.Question{ID: id, WorkspaceID: old.WorkspaceID, QuestionText: write.QuestionText, Kind: old.Kind, AnswerMarkdown: write.AnswerMarkdown, Status: write.Status, Priority: write.Priority, DueDate: write.DueDate}
	if err := domain.ValidateTransition(old.Status, q.Status, q.AnswerMarkdown); err != nil {
		return domain.Question{}, err
	}
	if err := q.Validate(); err != nil {
		return domain.Question{}, err
	}
	if write.WorkspaceID != "" && write.WorkspaceID != old.WorkspaceID {
		return domain.Question{}, ErrWorkspaceBoundary
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Question{}, err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, "UPDATE questions SET question_text=?,answer_markdown=?,status=?,priority=?,due_date=?,updated_at=?,version=version+1 WHERE id=? AND version=?", q.QuestionText, q.AnswerMarkdown, q.Status, q.Priority, q.DueDate, s.now(), id, version)
	if err != nil {
		return domain.Question{}, err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return domain.Question{}, ErrConflict
	}
	if err := s.replaceQuestionTagsTx(ctx, tx, id, old.WorkspaceID, write.TagIDs); err != nil {
		return domain.Question{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.Question{}, err
	}
	return s.GetQuestion(ctx, id)
}

func (s *CaptureStore) CreateNote(ctx context.Context, write domain.NoteWrite) (domain.Note, error) {
	if err := write.Validate(); err != nil {
		return domain.Note{}, err
	}
	ids, err := domain.ParseQuestionDirectives(write.BodyMarkdown)
	if err != nil {
		return domain.Note{}, err
	}
	if len(ids) != len(write.QuestionLinks) {
		return domain.Note{}, fmt.Errorf("question directives and links must match")
	}
	for i, id := range ids {
		if id != write.QuestionLinks[i].QuestionID {
			return domain.Note{}, fmt.Errorf("question link order does not match note directives")
		}
	}
	id, err := s.newID()
	if err != nil {
		return domain.Note{}, err
	}
	now := s.now()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Note{}, err
	}
	defer tx.Rollback()
	if err := s.checkNoteReferences(ctx, tx, "", write.WorkspaceID, write.TopicID, write.ParentNoteID); err != nil {
		return domain.Note{}, err
	}
	if _, err = tx.ExecContext(ctx, "INSERT INTO notes(id,workspace_id,topic_id,parent_note_id,title,body_markdown,created_at,updated_at,version) VALUES(?,?,?,?,?,?,?,?,1)", id, write.WorkspaceID, write.TopicID, write.ParentNoteID, strings.TrimSpace(write.Title), write.BodyMarkdown, now, now); err != nil {
		return domain.Note{}, err
	}
	for _, link := range write.QuestionLinks {
		if err := s.insertLink(ctx, tx, id, write.WorkspaceID, link); err != nil {
			return domain.Note{}, err
		}
	}
	if err := s.replaceNoteTagsTx(ctx, tx, id, write.WorkspaceID, write.TagIDs); err != nil {
		return domain.Note{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.Note{}, err
	}
	return s.GetNote(ctx, id)
}
func (s *CaptureStore) GetNote(ctx context.Context, id string) (domain.Note, error) {
	var n domain.Note
	var topic, parent sql.NullString
	var created, updated string
	err := s.db.QueryRowContext(ctx, "SELECT id,workspace_id,topic_id,parent_note_id,title,body_markdown,created_at,updated_at,version FROM notes WHERE id=?", id).Scan(&n.ID, &n.WorkspaceID, &topic, &parent, &n.Title, &n.BodyMarkdown, &created, &updated, &n.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return n, ErrNotFound
	}
	if err != nil {
		return n, err
	}
	if topic.Valid {
		n.TopicID = &topic.String
	}
	if parent.Valid {
		n.ParentNoteID = &parent.String
	}
	n.CreatedAt, err = platform.ParseTimestamp(created)
	if err != nil {
		return n, err
	}
	n.UpdatedAt, err = platform.ParseTimestamp(updated)
	if err != nil {
		return n, err
	}
	n.QuestionLinks, err = s.noteLinks(ctx, id)
	if err != nil {
		return n, err
	}
	n.TagIDs, err = s.noteTags(ctx, id)
	return n, err
}
func (s *CaptureStore) UpdateNote(ctx context.Context, id string, write domain.NoteWrite, version int64) (domain.Note, error) {
	if err := write.Validate(); err != nil {
		return domain.Note{}, err
	}
	ids, err := domain.ParseQuestionDirectives(write.BodyMarkdown)
	if err != nil {
		return domain.Note{}, err
	}
	if len(ids) != len(write.QuestionLinks) {
		return domain.Note{}, fmt.Errorf("question directives and links must match")
	}
	for i, qid := range ids {
		if qid != write.QuestionLinks[i].QuestionID {
			return domain.Note{}, fmt.Errorf("question link order does not match note directives")
		}
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Note{}, err
	}
	defer tx.Rollback()
	var workspace string
	var current int64
	if err := tx.QueryRowContext(ctx, "SELECT workspace_id,version FROM notes WHERE id=?", id).Scan(&workspace, &current); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Note{}, ErrNotFound
		}
		return domain.Note{}, err
	}
	if current != version {
		return domain.Note{}, ErrConflict
	}
	if workspace != write.WorkspaceID {
		return domain.Note{}, ErrWorkspaceBoundary
	}
	if err := s.checkNoteReferences(ctx, tx, id, workspace, write.TopicID, write.ParentNoteID); err != nil {
		return domain.Note{}, err
	}
	now := s.now()
	if _, err := tx.ExecContext(ctx, "UPDATE notes SET topic_id=?,parent_note_id=?,title=?,body_markdown=?,updated_at=?,version=version+1 WHERE id=? AND version=?", write.TopicID, write.ParentNoteID, strings.TrimSpace(write.Title), write.BodyMarkdown, now, id, version); err != nil {
		return domain.Note{}, err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM note_questions WHERE note_id=?", id); err != nil {
		return domain.Note{}, err
	}
	for _, link := range write.QuestionLinks {
		if err := s.insertLink(ctx, tx, id, workspace, link); err != nil {
			return domain.Note{}, err
		}
	}
	if err := s.replaceNoteTagsTx(ctx, tx, id, workspace, write.TagIDs); err != nil {
		return domain.Note{}, err
	}
	if err := tx.Commit(); err != nil {
		return domain.Note{}, err
	}
	return s.GetNote(ctx, id)
}
func (s *CaptureStore) ReplaceNoteQuestions(ctx context.Context, noteID string, links []domain.NoteQuestionWrite) error {
	note, err := s.GetNote(ctx, noteID)
	if err != nil {
		return err
	}
	bodyIDs, err := domain.ParseQuestionDirectives(note.BodyMarkdown)
	if err != nil {
		return err
	}
	if len(bodyIDs) != len(links) {
		return ErrInvalidLink
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, "DELETE FROM note_questions WHERE note_id=?", noteID); err != nil {
		return err
	}
	for i, link := range links {
		if i >= len(bodyIDs) || bodyIDs[i] != link.QuestionID {
			return ErrInvalidLink
		}
		if err := s.insertLink(ctx, tx, noteID, note.WorkspaceID, link); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *CaptureStore) checkNoteReferences(ctx context.Context, tx *sql.Tx, noteID, workspace string, topic, parent *string) error {
	var exists int
	if topic != nil {
		if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM topics WHERE id=? AND workspace_id=?", *topic, workspace).Scan(&exists); err != nil {
			return err
		}
		if exists != 1 {
			return ErrWorkspaceBoundary
		}
	}
	if parent != nil {
		if *parent == "" {
			return fmt.Errorf("parent note id is empty")
		}
		if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM notes WHERE id=? AND workspace_id=?", *parent, workspace).Scan(&exists); err != nil {
			return err
		}
		if exists != 1 {
			return ErrWorkspaceBoundary
		}
		candidate := *parent
		seen := map[string]bool{}
		for candidate != "" {
			if candidate == noteID || seen[candidate] {
				return fmt.Errorf("note hierarchy cycle")
			}
			seen[candidate] = true
			var next sql.NullString
			if err := tx.QueryRowContext(ctx, "SELECT parent_note_id FROM notes WHERE id=?", candidate).Scan(&next); err != nil {
				return err
			}
			if next.Valid {
				candidate = next.String
			} else {
				candidate = ""
			}
		}
	}
	return nil
}
func (s *CaptureStore) insertLink(ctx context.Context, tx *sql.Tx, noteID, workspace string, link domain.NoteQuestionWrite) error {
	if link.Position < 0 {
		return ErrInvalidLink
	}
	if err := domain.ValidateDisplayMode(link.DisplayMode); err != nil {
		return err
	}
	var qworkspace string
	if err := tx.QueryRowContext(ctx, "SELECT workspace_id FROM questions WHERE id=?", link.QuestionID).Scan(&qworkspace); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidLink
		}
		return err
	}
	if qworkspace != workspace {
		return ErrWorkspaceBoundary
	}
	_, err := tx.ExecContext(ctx, "INSERT INTO note_questions(note_id,question_id,display_mode,position,created_at,updated_at) VALUES(?,?,?,?,?,?)", noteID, link.QuestionID, link.DisplayMode, link.Position, s.now(), s.now())
	return err
}
func (s *CaptureStore) noteLinks(ctx context.Context, id string) ([]domain.NoteQuestion, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT nq.question_id,nq.display_mode,nq.position,q.workspace_id,q.question_text,q.kind,q.status,q.priority,q.due_date,q.version FROM note_questions nq JOIN questions q ON q.id=nq.question_id WHERE nq.note_id=? ORDER BY nq.position`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.NoteQuestion, 0)
	for rows.Next() {
		var l domain.NoteQuestion
		var due sql.NullString
		if err := rows.Scan(&l.QuestionID, &l.DisplayMode, &l.Position, &l.Question.WorkspaceID, &l.Question.QuestionText, &l.Question.Kind, &l.Question.Status, &l.Question.Priority, &due, &l.Question.Version); err != nil {
			return nil, err
		}
		l.Question.ID = l.QuestionID
		if due.Valid {
			l.Question.DueDate = &due.String
		}
		result = append(result, l)
	}
	return result, rows.Err()
}
func (s *CaptureStore) linkedNotes(ctx context.Context, id string) ([]domain.LinkedNoteSummary, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT n.id,n.title,nq.display_mode FROM note_questions nq JOIN notes n ON n.id=nq.note_id WHERE nq.question_id=? ORDER BY n.title COLLATE NOCASE,n.id", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.LinkedNoteSummary, 0)
	for rows.Next() {
		var value domain.LinkedNoteSummary
		if err := rows.Scan(&value.ID, &value.Title, &value.DisplayMode); err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	return result, rows.Err()
}
func (s *CaptureStore) questionTags(ctx context.Context, id string) ([]string, error) {
	return stringIDs(ctx, s.db, "SELECT tag_id FROM question_tags WHERE question_id=? ORDER BY tag_id", id)
}
func (s *CaptureStore) noteTags(ctx context.Context, id string) ([]string, error) {
	return stringIDs(ctx, s.db, "SELECT tag_id FROM note_tags WHERE note_id=? ORDER BY tag_id", id)
}
func stringIDs(ctx context.Context, db interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}, query string, args ...any) ([]string, error) {
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result = append(result, id)
	}
	return result, rows.Err()
}
func (s *CaptureStore) reminder(ctx context.Context, id string) (*domain.Reminder, error) {
	var r domain.Reminder
	var scheduled, created, updated string
	var last sql.NullString
	err := s.db.QueryRowContext(ctx, "SELECT id,scheduled_at,state,last_evaluated_at,created_at,updated_at,version FROM reminders WHERE question_id=?", id).Scan(&r.ID, &scheduled, &r.State, &last, &created, &updated, &r.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	r.ScheduledAt, err = platform.ParseTimestamp(scheduled)
	if err != nil {
		return nil, err
	}
	if last.Valid {
		value, e := platform.ParseTimestamp(last.String)
		if e != nil {
			return nil, e
		}
		r.LastEvaluatedAt = &value
	}
	r.CreatedAt, err = platform.ParseTimestamp(created)
	if err != nil {
		return nil, err
	}
	r.UpdatedAt, err = platform.ParseTimestamp(updated)
	return &r, err
}
func (s *CaptureStore) replaceQuestionTags(ctx context.Context, _ []string, id, workspace string, tags []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.replaceQuestionTagsTx(ctx, tx, id, workspace, tags); err != nil {
		return err
	}
	return tx.Commit()
}
func (s *CaptureStore) replaceQuestionTagsTx(ctx context.Context, tx *sql.Tx, id, workspace string, tags []string) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM question_tags WHERE question_id=?", id); err != nil {
		return err
	}
	for _, tag := range tags {
		var owner string
		if err := tx.QueryRowContext(ctx, "SELECT owner_workspace_id FROM tags WHERE id=?", tag).Scan(&owner); err != nil {
			return err
		}
		if owner != workspace {
			var available int
			if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM tag_workspace_access WHERE tag_id=? AND workspace_id=?", tag, workspace).Scan(&available); err != nil {
				return err
			}
			if available != 1 {
				return ErrWorkspaceBoundary
			}
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO question_tags(question_id,tag_id,created_at) VALUES(?,?,?)", id, tag, s.now()); err != nil {
			return err
		}
	}
	return nil
}
func (s *CaptureStore) replaceNoteTagsTx(ctx context.Context, tx *sql.Tx, id, workspace string, tags []string) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM note_tags WHERE note_id=?", id); err != nil {
		return err
	}
	for _, tag := range tags {
		var owner string
		if err := tx.QueryRowContext(ctx, "SELECT owner_workspace_id FROM tags WHERE id=?", tag).Scan(&owner); err != nil {
			return err
		}
		if owner != workspace {
			var available int
			if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM tag_workspace_access WHERE tag_id=? AND workspace_id=?", tag, workspace).Scan(&available); err != nil {
				return err
			}
			if available != 1 {
				return ErrWorkspaceBoundary
			}
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO note_tags(note_id,tag_id,created_at) VALUES(?,?,?)", id, tag, s.now()); err != nil {
			return err
		}
	}
	return nil
}
func isUnique(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint")
}

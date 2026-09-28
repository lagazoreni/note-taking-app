package sqlite

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"noted.local/noted/internal/domain"
)

type queryCursor struct {
	Value string `json:"v"`
	ID    string `json:"i"`
}

func EncodeCursor(value, id string) string {
	data, _ := json.Marshal(queryCursor{Value: value, ID: id})
	return base64.RawURLEncoding.EncodeToString(data)
}
func DecodeCursor(cursor string) (value, id string, err error) {
	if cursor == "" {
		return "", "", nil
	}
	data, e := base64.RawURLEncoding.DecodeString(cursor)
	if e != nil {
		return "", "", fmt.Errorf("invalid cursor")
	}
	var c queryCursor
	if e = json.Unmarshal(data, &c); e != nil || c.Value == "" || c.ID == "" {
		return "", "", fmt.Errorf("invalid cursor")
	}
	return c.Value, c.ID, nil
}

func (s *CaptureStore) QueryQuestions(ctx context.Context, input domain.QuestionQuery) (domain.QuestionPage, error) {
	query, err := domain.NormalizeQuestionQuery(input)
	if err != nil {
		return domain.QuestionPage{}, err
	}
	where := []string{"q.workspace_id = ?"}
	args := []any{query.WorkspaceID}
	if query.Kind != "" {
		where = append(where, "q.kind = ?")
		args = append(args, query.Kind)
	}
	if len(query.Statuses) > 0 {
		marks := make([]string, len(query.Statuses))
		for i, status := range query.Statuses {
			marks[i] = "?"
			args = append(args, status)
		}
		where = append(where, "q.status IN ("+strings.Join(marks, ",")+")")
	}
	if query.TopicID != nil {
		where = append(where, "EXISTS (SELECT 1 FROM note_questions nq_topic JOIN notes n_topic ON n_topic.id=nq_topic.note_id WHERE nq_topic.question_id=q.id AND n_topic.topic_id=?)")
		args = append(args, *query.TopicID)
	}
	if query.TagID != nil {
		where = append(where, "EXISTS (SELECT 1 FROM question_tags qt_filter WHERE qt_filter.question_id=q.id AND qt_filter.tag_id=?)")
		args = append(args, *query.TagID)
	}
	dateFilters := []struct {
		value  *string
		column string
	}{{query.CreatedFrom, "q.created_at >= ?"}, {query.CreatedTo, "q.created_at <= ?"}, {query.UpdatedFrom, "q.updated_at >= ?"}, {query.UpdatedTo, "q.updated_at <= ?"}, {query.DueFrom, "q.due_date >= ?"}, {query.DueTo, "q.due_date <= ?"}}
	for _, filter := range dateFilters {
		if filter.value != nil {
			where = append(where, filter.column)
			args = append(args, *filter.value)
		}
	}
	if query.HasAnswer != nil {
		if *query.HasAnswer {
			where = append(where, "q.answer_markdown IS NOT NULL AND length(trim(q.answer_markdown)) > 0")
		} else {
			where = append(where, "(q.answer_markdown IS NULL OR length(trim(q.answer_markdown)) = 0)")
		}
	}
	if query.HasLinkedNotes != nil {
		if *query.HasLinkedNotes {
			where = append(where, "EXISTS (SELECT 1 FROM note_questions nq_linked WHERE nq_linked.question_id=q.id")
		} else {
			where = append(where, "NOT EXISTS (SELECT 1 FROM note_questions nq_linked WHERE nq_linked.question_id=q.id)")
		}
	}
	column := map[string]string{"createdAt": "q.created_at", "updatedAt": "q.updated_at", "dueDate": "q.due_date", "priority": "q.priority", "questionText": "q.question_text", "status": "q.status"}[query.Sort]
	if column == "" {
		column = "q.updated_at"
	}
	direction := strings.ToUpper(query.Direction)
	comparison := "<"
	if direction == "ASC" {
		comparison = ">"
	}
	if query.Cursor != "" {
		value, id, e := DecodeCursor(query.Cursor)
		if e != nil {
			return domain.QuestionPage{}, e
		}
		where = append(where, "(("+column+" IS NOT NULL AND "+column+" "+comparison+" ?) OR ("+column+" = ? AND q.id "+comparison+" ?))")
		args = append(args, value, value, id)
	}
	order := column + " " + direction + ", q.id " + direction
	sqlQuery := "SELECT q.id FROM questions q WHERE " + strings.Join(where, " AND ") + " ORDER BY " + order + " LIMIT ?"
	args = append(args, query.PageSize+1)
	rows, e := s.db.QueryContext(ctx, sqlQuery, args...)
	if e != nil {
		return domain.QuestionPage{}, e
	}
	defer rows.Close()
	ids := make([]string, 0, query.PageSize+1)
	for rows.Next() {
		var id string
		if e := rows.Scan(&id); e != nil {
			return domain.QuestionPage{}, e
		}
		ids = append(ids, id)
	}
	if e := rows.Err(); e != nil {
		return domain.QuestionPage{}, e
	}
	next := (*string)(nil)
	if len(ids) > query.PageSize {
		last := ids[query.PageSize-1]
		ids = ids[:query.PageSize]
		q, e := s.getQuestionBase(ctx, last)
		if e != nil {
			return domain.QuestionPage{}, e
		}
		value := questionSortValue(q, query.Sort)
		encoded := EncodeCursor(value, q.ID)
		next = &encoded
	}
	items := make([]domain.Question, 0, len(ids))
	for _, id := range ids {
		q, e := s.GetQuestion(ctx, id)
		if e != nil {
			return domain.QuestionPage{}, e
		}
		items = append(items, q)
	}
	return domain.QuestionPage{Items: items, NextCursor: next}, nil
}
func questionSortValue(q domain.Question, sort string) string {
	switch sort {
	case "createdAt":
		return q.CreatedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00")
	case "dueDate":
		if q.DueDate != nil {
			return *q.DueDate
		}
		return ""
	case "priority":
		return string(q.Priority)
	case "questionText":
		return q.QuestionText
	case "status":
		return string(q.Status)
	default:
		return q.UpdatedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00")
	}
}

func (s *CaptureStore) ListQuestions(ctx context.Context, workspaceID string, statuses []domain.QuestionStatus) ([]domain.Question, error) {
	query := "SELECT id FROM questions WHERE workspace_id = ?"
	args := []any{workspaceID}
	if len(statuses) > 0 {
		marks := make([]string, len(statuses))
		for i, status := range statuses {
			if !domain.ValidStatus(status) {
				return nil, fmt.Errorf("invalid question status")
			}
			marks[i] = "?"
			args = append(args, status)
		}
		query += " AND status IN (" + strings.Join(marks, ",") + ")"
	}
	query += " ORDER BY updated_at DESC, id DESC"
	rows, err := s.db.QueryContext(ctx, query, args...)
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

func (s *CaptureStore) Count(ctx context.Context, table string) (int64, error) {
	allowed := map[string]bool{"workspaces": true, "notes": true, "questions": true, "note_questions": true}
	if !allowed[table] {
		return 0, fmt.Errorf("unsupported table")
	}
	var count int64
	err := s.db.QueryRowContext(ctx, "SELECT count(*) FROM "+table).Scan(&count)
	return count, err
}

func scanNullableString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	return &value.String
}

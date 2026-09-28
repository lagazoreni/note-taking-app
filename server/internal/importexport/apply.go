package importexport

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

var importOrder = []struct{ path, table string }{{"data/workspaces.json", "workspaces"}, {"data/topics.json", "topics"}, {"data/notes.json", "notes"}, {"data/questions.json", "questions"}, {"data/tags.json", "tags"}, {"data/tag-workspace-access.json", "tag_workspace_access"}, {"data/note-tags.json", "note_tags"}, {"data/question-tags.json", "question_tags"}, {"data/note-questions.json", "note_questions"}, {"data/reminders.json", "reminders"}}

func Apply(ctx context.Context, db *sql.DB, archive Archive, resolutions map[string]string) (map[string]int, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	counts := map[string]int{}
	for _, entry := range importOrder {
		data := archive.Files[entry.path]
		var records []map[string]any
		if err := json.Unmarshal(data, &records); err != nil {
			return nil, err
		}
		for _, record := range records {
			id := fmt.Sprint(record["id"])
			if id == "<nil>" {
				id = fmt.Sprint(record["noteId"])
			}
			if resolution, ok := resolutions[id]; ok && resolution == "keep_existing" {
				continue
			}
			if err := insertRecord(ctx, tx, entry.table, record); err != nil {
				return nil, fmt.Errorf("import %s: %w", entry.table, err)
			}
			counts[entry.table]++
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return counts, nil
}
func insertRecord(ctx context.Context, tx *sql.Tx, table string, record map[string]any) error {
	allowed := map[string]map[string]bool{"workspaces": {"id": true, "name": true, "createdAt": true, "updatedAt": true, "version": true}, "topics": {"id": true, "workspaceId": true, "name": true, "createdAt": true, "updatedAt": true, "version": true}, "notes": {"id": true, "workspaceId": true, "topicId": true, "parentNoteId": true, "title": true, "bodyMarkdown": true, "createdAt": true, "updatedAt": true, "version": true}, "questions": {"id": true, "workspaceId": true, "questionText": true, "kind": true, "answerMarkdown": true, "status": true, "priority": true, "dueDate": true, "createdAt": true, "updatedAt": true, "version": true}, "tags": {"id": true, "ownerWorkspaceId": true, "name": true, "createdAt": true, "updatedAt": true, "version": true}, "tag_workspace_access": {"tagId": true, "workspaceId": true, "createdAt": true}, "note_tags": {"noteId": true, "tagId": true, "createdAt": true}, "question_tags": {"questionId": true, "tagId": true, "createdAt": true}, "note_questions": {"noteId": true, "questionId": true, "displayMode": true, "position": true, "createdAt": true, "updatedAt": true}, "reminders": {"id": true, "questionId": true, "scheduledAt": true, "state": true, "lastEvaluatedAt": true, "createdAt": true, "updatedAt": true, "version": true}}
	fields := allowed[table]
	if fields == nil {
		return fmt.Errorf("unsupported table")
	}
	columns := make([]string, 0)
	values := make([]any, 0)
	for key, value := range record {
		if !fields[key] {
			continue
		}
		columns = append(columns, toSnake(key))
		values = append(values, value)
	}
	if len(columns) == 0 {
		return nil
	}
	marks := make([]string, len(columns))
	for i := range marks {
		marks[i] = "?"
	}
	query := "INSERT OR REPLACE INTO " + table + " (" + strings.Join(columns, ",") + ") VALUES (" + strings.Join(marks, ",") + ")"
	_, err := tx.ExecContext(ctx, query, values...)
	return err
}
func toSnake(value string) string {
	var out strings.Builder
	for i, r := range value {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				out.WriteByte('_')
			}
			out.WriteRune(r + 'a' - 'A')
		} else {
			out.WriteRune(r)
		}
	}
	return out.String()
}

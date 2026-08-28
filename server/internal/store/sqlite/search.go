package sqlite

import (
	"context"
	"fmt"
	"html"
	"regexp"
	"strings"

	"noted.local/noted/internal/domain"
)

var searchWord = regexp.MustCompile(`[[:alnum:]]+`)

func ftsQuery(input string) string {
	words := searchWord.FindAllString(strings.ToLower(input), -1)
	parts := make([]string, 0, len(words))
	for _, word := range words {
		parts = append(parts, `"`+strings.ReplaceAll(word, `"`, `""`)+`"*`)
	}
	return strings.Join(parts, " OR ")
}

func (s *CaptureStore) Search(ctx context.Context, input string, contentScope domain.SearchContentScope, workspaceScope domain.SearchWorkspaceScope, workspaceID *string, limit int) ([]domain.SearchResult, error) {
	if strings.TrimSpace(input) == "" {
		return []domain.SearchResult{}, nil
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	match := ftsQuery(input)
	if match == "" {
		return []domain.SearchResult{}, nil
	}
	if workspaceScope == domain.SearchCurrent && (workspaceID == nil || *workspaceID == "") {
		return nil, fmt.Errorf("workspace is required for current scope")
	}
	result := make([]domain.SearchResult, 0)
	includeNotes := contentScope == domain.SearchNotes || contentScope == domain.SearchEverything
	includeQuestions := contentScope == domain.SearchQuestions || contentScope == domain.SearchEverything
	includeAnswers := contentScope == domain.SearchAnswers || contentScope == domain.SearchEverything
	if includeNotes {
		items, err := s.searchFTS(ctx, "notes", match, workspaceScope, workspaceID, limit)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
	}
	if includeQuestions {
		items, err := s.searchFTS(ctx, "questions", match, workspaceScope, workspaceID, limit)
		if err != nil {
			return nil, err
		}
		for i := range items {
			items[i].Type = "question"
		}
		result = append(result, items...)
	}
	if includeAnswers {
		items, err := s.searchFTSAnswers(ctx, match, workspaceScope, workspaceID, limit)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
	}
	if len(result) > limit {
		result = result[:limit]
	}
	return result, nil
}
func (s *CaptureStore) searchFTS(ctx context.Context, kind, match string, scope domain.SearchWorkspaceScope, workspaceID *string, limit int) ([]domain.SearchResult, error) {
	table := kind + "_fts"
	idColumn := "note_id"
	if kind == "questions" {
		idColumn = "question_id"
	}
	where := ""
	args := []any{match}
	if scope == domain.SearchCurrent {
		where = " AND f.workspace_id=?"
		args = append(args, *workspaceID)
	}
	args = append(args, limit)
	query := fmt.Sprintf(`SELECT f.%s,w.id,w.name,CASE WHEN ?='notes' THEN n.title ELSE '' END,snippet(f,CASE WHEN ?='notes' THEN 2 ELSE 2 END,'<mark>','</mark>','…',24),bm25(f) FROM %s f JOIN %s entity ON entity.id=f.%s JOIN workspaces w ON w.id=entity.workspace_id LEFT JOIN notes n ON n.id=f.note_id WHERE %s MATCH ? %s ORDER BY bm25(f),f.%s LIMIT ?`, idColumn, kind, kind, idColumn, table, where, idColumn)
	// The query above needs the kind literal for CASE but the FTS MATCH parameter is first in the SQL in practice; use a simpler explicit query per table.
	if kind == "notes" {
		query = `SELECT f.note_id,w.id,w.name,n.title,snippet(notes_fts,3,'<mark>','</mark>','…',24),bm25(notes_fts) FROM notes_fts f JOIN notes n ON n.id=f.note_id JOIN workspaces w ON w.id=n.workspace_id WHERE notes_fts MATCH ?` + where + ` ORDER BY bm25(notes_fts),f.note_id LIMIT ?`
	} else {
		query = `SELECT f.question_id,w.id,w.name,q.question_text,snippet(questions_fts,2,'<mark>','</mark>','…',24),bm25(questions_fts) FROM questions_fts f JOIN questions q ON q.id=f.question_id JOIN workspaces w ON w.id=q.workspace_id WHERE questions_fts MATCH ?` + where + ` ORDER BY bm25(questions_fts),f.question_id LIMIT ?`
	}
	args = []any{match}
	if scope == domain.SearchCurrent {
		args = append(args, *workspaceID)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.SearchResult, 0)
	for rows.Next() {
		var value domain.SearchResult
		var snippet string
		var rank float64
		if err := rows.Scan(&value.ID, &value.WorkspaceID, &value.WorkspaceName, &value.Title, &snippet, &rank); err != nil {
			return nil, err
		}
		value.Type = map[string]string{"notes": "note", "questions": "question"}[kind]
		value.Snippet = html.EscapeString(strings.ReplaceAll(strings.ReplaceAll(snippet, "<mark>", ""), "</mark>", ""))
		if value.Type == "note" {
			value.Destination = "/notes/" + value.ID
		} else {
			value.Destination = "/questions/" + value.ID
		}
		result = append(result, value)
	}
	return result, rows.Err()
}
func (s *CaptureStore) searchFTSAnswers(ctx context.Context, match string, scope domain.SearchWorkspaceScope, workspaceID *string, limit int) ([]domain.SearchResult, error) {
	where := ""
	args := []any{match}
	if scope == domain.SearchCurrent {
		where = " AND f.workspace_id=?"
		args = append(args, *workspaceID)
	}
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, `SELECT f.question_id,w.id,w.name,q.question_text,snippet(questions_fts,3,'<mark>','</mark>','…',24),bm25(questions_fts) FROM questions_fts f JOIN questions q ON q.id=f.question_id JOIN workspaces w ON w.id=q.workspace_id WHERE questions_fts MATCH ? AND length(trim(COALESCE(f.answer_markdown,''))) > 0`+where+` ORDER BY bm25(questions_fts),f.question_id LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.SearchResult, 0)
	for rows.Next() {
		var value domain.SearchResult
		var snippet string
		var rank float64
		if err := rows.Scan(&value.ID, &value.WorkspaceID, &value.WorkspaceName, &value.Title, &snippet, &rank); err != nil {
			return nil, err
		}
		value.Type = "answer"
		value.Snippet = html.EscapeString(strings.ReplaceAll(strings.ReplaceAll(snippet, "<mark>", ""), "</mark>", ""))
		value.Destination = "/questions/" + value.ID
		result = append(result, value)
	}
	return result, rows.Err()
}

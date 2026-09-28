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

// Search keeps the original call shape compatible with callers that only pass
// a limit. New callers may pass a *string (or string) tag ID followed by the
// limit. SearchWithTag is the typed entry point used by the HTTP handler.
func (s *CaptureStore) Search(ctx context.Context, input string, contentScope domain.SearchContentScope, workspaceScope domain.SearchWorkspaceScope, workspaceID *string, options ...any) ([]domain.SearchResult, error) {
	tagID, limit := searchOptions(options)
	return s.search(ctx, input, contentScope, workspaceScope, workspaceID, tagID, limit)
}

func (s *CaptureStore) SearchWithTag(ctx context.Context, input string, contentScope domain.SearchContentScope, workspaceScope domain.SearchWorkspaceScope, workspaceID, tagID *string, limit int) ([]domain.SearchResult, error) {
	return s.search(ctx, input, contentScope, workspaceScope, workspaceID, tagID, limit)
}

func searchOptions(options []any) (*string, int) {
	limit := 50
	var tagID *string
	for _, option := range options {
		switch value := option.(type) {
		case int:
			limit = value
		case *string:
			tagID = value
		case string:
			if value != "" {
				copy := value
				tagID = &copy
			}
		}
	}
	return tagID, limit
}

func (s *CaptureStore) search(ctx context.Context, input string, contentScope domain.SearchContentScope, workspaceScope domain.SearchWorkspaceScope, workspaceID, tagID *string, limit int) ([]domain.SearchResult, error) {
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
	if tagID != nil && strings.TrimSpace(*tagID) == "" {
		tagID = nil
	}
	result := make([]domain.SearchResult, 0)
	includeNotes := contentScope == domain.SearchNotes || contentScope == domain.SearchEverything
	includeQuestions := contentScope == domain.SearchQuestions || contentScope == domain.SearchEverything
	includeAnswers := contentScope == domain.SearchAnswers || contentScope == domain.SearchEverything
	if includeNotes {
		items, err := s.searchFTS(ctx, "notes", match, workspaceScope, workspaceID, tagID, limit)
		if err != nil {
			return nil, err
		}
		result = append(result, items...)
	}
	if includeQuestions {
		items, err := s.searchFTS(ctx, "questions", match, workspaceScope, workspaceID, nil, limit)
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

func (s *CaptureStore) searchFTS(ctx context.Context, kind, match string, scope domain.SearchWorkspaceScope, workspaceID, tagID *string, limit int) ([]domain.SearchResult, error) {
	where := ""
	args := []any{match}
	if scope == domain.SearchCurrent {
		where = " AND f.workspace_id=?"
		args = append(args, *workspaceID)
	}
	if kind == "notes" && tagID != nil {
		// EXISTS avoids multiplying a note when a tag has more than one
		// relationship and can use note_tags_tag_note(tag_id, note_id).
		where += " AND EXISTS (SELECT 1 FROM note_tags nt WHERE nt.note_id=f.note_id AND nt.tag_id=?)"
		args = append(args, *tagID)
	}
	args = append(args, limit)

	var query string
	if kind == "notes" {
		query = `SELECT f.note_id,w.id,w.name,n.title,snippet(notes_fts,3,'<mark>','</mark>','…',24),bm25(notes_fts) FROM notes_fts f JOIN notes n ON n.id=f.note_id JOIN workspaces w ON w.id=n.workspace_id WHERE notes_fts MATCH ?` + where + ` ORDER BY bm25(notes_fts),f.note_id LIMIT ?`
	} else {
		query = `SELECT f.question_id,w.id,w.name,q.question_text,snippet(questions_fts,2,'<mark>','</mark>','…',24),bm25(questions_fts) FROM questions_fts f JOIN questions q ON q.id=f.question_id JOIN workspaces w ON w.id=q.workspace_id WHERE questions_fts MATCH ?` + where + ` ORDER BY bm25(questions_fts),f.question_id LIMIT ?`
	}

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

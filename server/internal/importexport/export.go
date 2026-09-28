package importexport

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"noted.local/noted/internal/platform"
)

var exportFiles = []struct{ path, entity, table, order string }{
	{"data/workspaces.json", "workspace", "workspaces", "id"},
	{"data/topics.json", "topic", "topics", "id"},
	{"data/notes.json", "note", "notes", "id"},
	{"data/questions.json", "question", "questions", "id"},
	{"data/note-questions.json", "noteQuestion", "note_questions", "note_id,question_id"},
	{"data/tags.json", "tag", "tags", "id"},
	{"data/tag-workspace-access.json", "tagWorkspaceAccess", "tag_workspace_access", "tag_id,workspace_id"},
	{"data/note-tags.json", "noteTag", "note_tags", "note_id,tag_id"},
	{"data/question-tags.json", "questionTag", "question_tags", "question_id,tag_id"},
	{"data/reminders.json", "reminder", "reminders", "id"},
}

func Export(ctx context.Context, db *sql.DB, applicationVersion string, clock platform.Clock) ([]byte, error) {
	if clock == nil {
		clock = platform.SystemClock{}
	}
	files := make(map[string][]byte, len(exportFiles))
	entries := make([]FileEntry, 0, len(exportFiles))
	counts := Counts{}
	for _, file := range exportFiles {
		data, count, err := queryRecords(ctx, db, file.table, file.order)
		if err != nil {
			return nil, fmt.Errorf("export %s: %w", file.table, err)
		}
		files[file.path] = data
		entries = append(entries, FileEntry{Path: file.path, Entity: file.entity, SHA256: sum(data), Records: count})
		setCount(&counts, file.entity, count)
	}
	manifest := Manifest{Format: "interactive-note-questions-export", FormatVersion: 1, ExportedAt: clock.Now().UTC(), ApplicationVersion: applicationVersion, Files: entries, Counts: counts}
	manifestBytes, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return nil, err
	}
	manifestBytes = append(manifestBytes, '\n')
	files["manifest.json"] = manifestBytes
	paths := make([]string, 0, len(files))
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	checksums := make([]string, 0, len(paths))
	for _, path := range paths {
		if path == "checksums.txt" {
			continue
		}
		checksums = append(checksums, sum(files[path])+"  "+path)
	}
	files["checksums.txt"] = []byte(strings.Join(checksums, "\n") + "\n")
	paths = paths[:0]
	for path := range files {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	var buffer bytes.Buffer
	writer := zip.NewWriter(&buffer)
	for _, path := range paths {
		header := &zip.FileHeader{Name: path, Method: zip.Deflate}
		header.SetModTime(time.Unix(0, 0))
		part, err := writer.CreateHeader(header)
		if err != nil {
			return nil, err
		}
		if _, err := part.Write(files[path]); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func queryRecords(ctx context.Context, db *sql.DB, table, order string) ([]byte, int, error) {
	allowed := map[string]bool{"workspaces": true, "topics": true, "notes": true, "questions": true, "note_questions": true, "tags": true, "tag_workspace_access": true, "note_tags": true, "question_tags": true, "reminders": true}
	if !allowed[table] {
		return nil, 0, fmt.Errorf("unsupported table")
	}
	rows, err := db.QueryContext(ctx, "SELECT * FROM "+table+" ORDER BY "+order)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, 0, err
	}
	records := make([]map[string]any, 0)
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return nil, 0, err
		}
		record := map[string]any{}
		for i, column := range columns {
			value := values[i]
			if data, ok := value.([]byte); ok {
				value = string(data)
			}
			record[toCamel(column)] = value
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return nil, 0, err
	}
	return append(data, '\n'), len(records), nil
}
func toCamel(value string) string {
	parts := strings.Split(value, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
func sum(data []byte) string { hash := sha256.Sum256(data); return hex.EncodeToString(hash[:]) }
func setCount(counts *Counts, entity string, value int) {
	switch entity {
	case "workspace":
		counts.Workspaces = value
	case "topic":
		counts.Topics = value
	case "note":
		counts.Notes = value
	case "question":
		counts.Questions = value
	case "noteQuestion":
		counts.NoteQuestions = value
	case "tag":
		counts.Tags = value
	case "tagWorkspaceAccess":
		counts.TagWorkspaceAccess = value
	case "noteTag":
		counts.NoteTags = value
	case "questionTag":
		counts.QuestionTags = value
	case "reminder":
		counts.Reminders = value
	}
}

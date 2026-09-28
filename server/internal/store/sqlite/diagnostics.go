package sqlite

import (
	"context"
	"database/sql"
	"fmt"
)

func IntegrityCheck(ctx context.Context, db *sql.DB) error {
	var result string
	if err := db.QueryRowContext(ctx, "PRAGMA integrity_check").Scan(&result); err != nil {
		return err
	}
	if result != "ok" {
		return fmt.Errorf("integrity check: %s", result)
	}
	return nil
}
func ForeignKeyCheck(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, "PRAGMA foreign_key_check")
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return fmt.Errorf("foreign-key violation")
	}
	return rows.Err()
}
func FTSCounts(ctx context.Context, db *sql.DB) (notes, questions int64, err error) {
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM notes_fts").Scan(&notes); err != nil {
		return 0, 0, err
	}
	if err := db.QueryRowContext(ctx, "SELECT count(*) FROM questions_fts").Scan(&questions); err != nil {
		return 0, 0, err
	}
	return notes, questions, nil
}
func RebuildFTS(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, "DELETE FROM notes_fts"); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, "INSERT INTO notes_fts(note_id,workspace_id,title,body_markdown) SELECT id,workspace_id,title,body_markdown FROM notes"); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, "DELETE FROM questions_fts"); err != nil {
		return err
	}
	_, err := db.ExecContext(ctx, "INSERT INTO questions_fts(question_id,workspace_id,question_text,answer_markdown) SELECT id,workspace_id,question_text,COALESCE(answer_markdown,'') FROM questions")
	return err
}

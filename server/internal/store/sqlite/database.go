package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const (
	maxOpenConnections = 4
	busyTimeoutMS      = 5000
)

type Database struct {
	*sql.DB
	Path string
}

func OpenDatabase(ctx context.Context, path string) (*Database, error) {
	db, err := OpenContext(ctx, path)
	if err != nil {
		return nil, err
	}
	return &Database{DB: db, Path: path}, nil
}

func Open(path string) (*sql.DB, error) {
	return OpenContext(context.Background(), path)
}

func OpenContext(ctx context.Context, path string) (*sql.DB, error) {
	if path == "" {
		return nil, fmt.Errorf("database path must not be empty")
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
	}
	dsn := "file:" + filepath.ToSlash(path) + fmt.Sprintf("?_pragma=foreign_keys(1)&_pragma=busy_timeout(%d)&_pragma=journal_mode(WAL)", busyTimeoutMS)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxOpenConnections)
	db.SetConnMaxIdleTime(0)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	for _, pragma := range []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		fmt.Sprintf("PRAGMA busy_timeout = %d", busyTimeoutMS),
		"PRAGMA synchronous = NORMAL",
	} {
		if _, err := db.ExecContext(ctx, pragma); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("configure sqlite (%s): %w", pragma, err)
		}
	}
	if err := CheckFTS5(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func CheckFTS5(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, "CREATE VIRTUAL TABLE temp.noted_fts5_check USING fts5(value)"); err != nil {
		return fmt.Errorf("sqlite FTS5 is unavailable: %w", err)
	}
	if _, err := db.ExecContext(ctx, "DROP TABLE temp.noted_fts5_check"); err != nil {
		return fmt.Errorf("remove FTS5 capability check: %w", err)
	}
	return nil
}

func Optimize(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "PRAGMA optimize")
	return err
}

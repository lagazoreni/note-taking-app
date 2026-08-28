package sqlite

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"noted.local/noted/internal/platform"
	"noted.local/noted/internal/store/migrations"
)

type Migration struct {
	Version  int
	Name     string
	SQL      []byte
	Checksum string
}

func LoadMigrations(files fs.FS) ([]Migration, error) {
	entries, err := fs.Glob(files, "*.sql")
	if err != nil {
		return nil, fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(entries)
	result := make([]Migration, 0, len(entries))
	for _, name := range entries {
		base := filepath.Base(name)
		parts := strings.SplitN(base, "_", 2)
		if len(parts) != 2 || !strings.HasSuffix(parts[1], ".sql") {
			return nil, fmt.Errorf("migration filename %q must be NNN_name.sql", name)
		}
		version, err := strconv.Atoi(parts[0])
		if err != nil || version <= 0 {
			return nil, fmt.Errorf("migration filename %q has invalid version", name)
		}
		contents, err := fs.ReadFile(files, name)
		if err != nil {
			return nil, fmt.Errorf("read migration %q: %w", name, err)
		}
		result = append(result, Migration{
			Version: version, Name: strings.TrimSuffix(parts[1], ".sql"), SQL: contents,
			Checksum: fmt.Sprintf("%x", sha256.Sum256(contents)),
		})
	}
	return result, nil
}

func RunMigrations(ctx context.Context, db *sql.DB) error {
	return RunMigrationsWithFS(ctx, db, migrations.Files, platform.SystemClock{})
}

func RunMigrationsWithFS(ctx context.Context, db *sql.DB, files embed.FS, clock platform.Clock) error {
	return runMigrations(ctx, db, files, clock)
}

func runMigrations(ctx context.Context, db *sql.DB, files fs.FS, clock platform.Clock) error {
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		checksum TEXT NOT NULL,
		applied_at TEXT NOT NULL
	)`); err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}
	all, err := LoadMigrations(files)
	if err != nil {
		return err
	}
	for _, migration := range all {
		var existing string
		err := db.QueryRowContext(ctx, "SELECT checksum FROM schema_migrations WHERE version = ?", migration.Version).Scan(&existing)
		switch err {
		case nil:
			if existing != migration.Checksum {
				return fmt.Errorf("migration %03d checksum mismatch", migration.Version)
			}
			continue
		case sql.ErrNoRows:
		default:
			return fmt.Errorf("check migration %03d: %w", migration.Version, err)
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("begin migration %03d: %w", migration.Version, err)
		}
		if _, err := tx.ExecContext(ctx, string(migration.SQL)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %03d: %w", migration.Version, err)
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations(version, name, checksum, applied_at) VALUES (?, ?, ?, ?)", migration.Version, migration.Name, migration.Checksum, platform.FormatTimestamp(clock.Now())); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %03d: %w", migration.Version, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %03d: %w", migration.Version, err)
		}
	}
	return nil
}

type MigrationRunner struct {
	Files embed.FS
	Clock platform.Clock
	_     time.Time
}

func (r MigrationRunner) Run(ctx context.Context, db *sql.DB) error {
	clock := r.Clock
	if clock == nil {
		clock = platform.SystemClock{}
	}
	return runMigrations(ctx, db, r.Files, clock)
}

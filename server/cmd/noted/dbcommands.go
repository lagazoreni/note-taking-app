package main

import (
	"context"
	"database/sql"

	"noted.local/noted/internal/store/sqlite"
)

func checkDatabase(ctx context.Context, db *sql.DB) error {
	if err := sqlite.IntegrityCheck(ctx, db); err != nil {
		return err
	}
	if err := sqlite.ForeignKeyCheck(ctx, db); err != nil {
		return err
	}
	return sqlite.Optimize(ctx, db)
}

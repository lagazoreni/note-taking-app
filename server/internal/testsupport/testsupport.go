package testsupport

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

func OpenTempDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	dir := t.TempDir()
	db, err := sql.Open("sqlite", "file:"+dir+"/test.db?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")
	if err != nil {
		t.Fatal(err)
	}
	db.SetMaxOpenConns(4)
	if err := db.PingContext(context.Background()); err != nil {
		db.Close()
		t.Fatal(err)
	}
	return db, func() { _ = db.Close() }
}

func HTTPServer(handler http.Handler) *httptest.Server { return httptest.NewServer(handler) }

func AssertForeignKeys(t *testing.T, db *sql.DB) {
	t.Helper()
	rows, err := db.Query("PRAGMA foreign_key_check")
	if err != nil {
		t.Fatalf("foreign key check: %v", err)
	}
	defer rows.Close()
	if rows.Next() {
		t.Fatal("foreign key check returned a violation")
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
}

func MustTempFile(t *testing.T, content []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "fixture-*")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	return f.Name()
}

// Package testutil provides shared helpers for integration and end-to-end
// tests that need a real Postgres connection. It intentionally imports
// "testing" (like net/http/httptest does) since it is only ever linked
// into test binaries.
package testutil

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // for mocking postgres db tag handling

	"github.com/eedriz99/go_blog/internal/env"
)

// OpenTestDB connects to the database used for local development
// (DB_ADDR, matching cmd/api/main.go's default) and skips the calling
// test if it isn't reachable, so `go test -tags=integration ./...` fails
// loudly only when it's actually supposed to run, not on machines without
// Postgres. The connection is closed automatically via t.Cleanup.
func OpenTestDB(t *testing.T) *sql.DB {
	t.Helper()

	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/go_blog?sslmode=disable")

	db, err := sql.Open("postgres", addr)
	if err != nil {
		t.Fatalf("open db %q: %v", addr, err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Skipf("skipping: postgres not reachable at %q: %v", addr, err)
	}

	t.Cleanup(func() { db.Close() })
	return db
}

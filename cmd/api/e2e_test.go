//go:build e2e

package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eedriz99/go_blog/internal/store"
	"github.com/eedriz99/go_blog/internal/testutil"
)

func newE2EServer(t *testing.T) (*httptest.Server, *sql.DB) {
	t.Helper()

	db := testutil.OpenTestDB(t)
	app := &application{
		config: config{env: "test", mail: mailConfig{expiry: time.Hour}},
		store:  store.NewStore(db),
	}

	srv := httptest.NewServer(app.mount())
	t.Cleanup(srv.Close)

	return srv, db
}

func postJSON(t *testing.T, url string, body any) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatalf("POST %s: %v", url, err)
	}
	return resp
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	return string(b)
}

// TestE2E_RegisterActivateLogin drives the full register -> activate ->
// login flow through a real HTTP server backed by a real Postgres
// database, exactly as a real client would. The API never hands the
// invitation token back over HTTP (it's meant to be emailed), so the test
// reads it straight out of the database, the same way a mail worker
// would source it in production.
func TestE2E_RegisterActivateLogin(t *testing.T) {
	srv, db := newE2EServer(t)

	suffix := fmt.Sprintf("%d", time.Now().UnixNano())
	email := fmt.Sprintf("e2e-%s@example.com", suffix)
	username := fmt.Sprintf("e2e-%s", suffix)
	password := "correct-horse-battery-staple"

	t.Cleanup(func() {
		if _, err := db.Exec(`DELETE FROM users WHERE email = $1`, email); err != nil {
			t.Logf("cleanup: failed to delete user %s: %v", email, err)
		}
	})

	registerBody := map[string]string{
		"email":      email,
		"first_name": "E2E",
		"last_name":  "Tester",
		"username":   username,
		"password":   password,
	}

	// 1. Register a new account.
	resp := postJSON(t, srv.URL+"/v1/auth/register", registerBody)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register: expected 201, got %d: %s", resp.StatusCode, readBody(t, resp))
	}
	resp.Body.Close()

	// 2. Registering the same email again is rejected as a conflict.
	resp = postJSON(t, srv.URL+"/v1/auth/register", registerBody)
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate register: expected 409, got %d: %s", resp.StatusCode, readBody(t, resp))
	}
	resp.Body.Close()

	// 3. The account isn't active yet, so logging in must fail even with
	// the correct password.
	resp = postJSON(t, srv.URL+"/v1/auth/login", map[string]string{"email": email, "password": password})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("login before activation: expected 401, got %d: %s", resp.StatusCode, readBody(t, resp))
	}
	resp.Body.Close()

	// 4. Pull the invitation token straight from the DB, standing in for
	// the activation email that would deliver it in production.
	var token string
	err := db.QueryRow(`
		SELECT ui.token FROM user_invitations ui
		JOIN users u ON u.id = ui.user_id
		WHERE u.email = $1`, email).Scan(&token)
	if err != nil {
		t.Fatalf("read invitation token: %v", err)
	}

	// 5. Activating with that token succeeds.
	resp, err = http.Post(srv.URL+"/v1/auth/activate/"+token, "application/json", nil)
	if err != nil {
		t.Fatalf("activate: %v", err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("activate: expected 202, got %d: %s", resp.StatusCode, readBody(t, resp))
	}
	resp.Body.Close()

	// 6. The token is single-use: activating again must not succeed a
	// second time.
	resp, err = http.Post(srv.URL+"/v1/auth/activate/"+token, "application/json", nil)
	if err != nil {
		t.Fatalf("re-activate: %v", err)
	}
	if resp.StatusCode == http.StatusAccepted {
		t.Fatalf("expected re-using a consumed activation token to fail, got 202")
	}
	resp.Body.Close()

	// 7. Wrong password is still rejected post-activation.
	resp = postJSON(t, srv.URL+"/v1/auth/login", map[string]string{"email": email, "password": "not-the-password"})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong password login: expected 401, got %d: %s", resp.StatusCode, readBody(t, resp))
	}
	resp.Body.Close()

	// 8. Correct credentials on an activated account succeed.
	resp = postJSON(t, srv.URL+"/v1/auth/login", map[string]string{"email": email, "password": password})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("login after activation: expected 200, got %d: %s", resp.StatusCode, readBody(t, resp))
	}
	resp.Body.Close()
}

// TestE2E_Register_RejectsUnknownFields exercises the strict JSON
// decoding (DisallowUnknownFields) end-to-end, guarding the contract that
// registration requests reject unexpected payload shapes rather than
// silently ignoring them.
func TestE2E_Register_RejectsUnknownFields(t *testing.T) {
	srv, _ := newE2EServer(t)

	resp := postJSON(t, srv.URL+"/v1/auth/register", map[string]any{
		"email":      fmt.Sprintf("e2e-bad-%d@example.com", time.Now().UnixNano()),
		"first_name": "E2E",
		"last_name":  "Tester",
		"username":   fmt.Sprintf("e2e-bad-%d", time.Now().UnixNano()),
		"password":   "whatever-password",
		"is_admin":   true, // not a recognized field
	})
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for an unknown field, got %d: %s", resp.StatusCode, readBody(t, resp))
	}
}

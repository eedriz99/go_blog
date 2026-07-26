package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eedriz99/go_blog/internal/dto/payload"
	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// mockUserStore implements the store.Storage.Users interface with
// swappable function fields, so each test wires up only the behavior it
// needs and nothing else.
type mockUserStore struct {
	createWithInvitationFn      func(ctx context.Context, u *model.User, expiry time.Duration) (*string, error)
	getByIDFn                   func(ctx context.Context, id string) (*model.User, error)
	getByEmailFn                func(ctx context.Context, email string) (*model.User, error)
	activateByInvitationTokenFn func(ctx context.Context, token string) error
	updateFn                    func(ctx context.Context, p payload.UpdateUserPayload) error
	deleteFn                    func(ctx context.Context, id string) error
	deleteTokenFn               func(ctx context.Context, token string) error
}

func (m *mockUserStore) CreateWithInvitation(ctx context.Context, u *model.User, expiry time.Duration) (*string, error) {
	return m.createWithInvitationFn(ctx, u, expiry)
}

func (m *mockUserStore) GetByID(ctx context.Context, id string) (*model.User, error) {
	return m.getByIDFn(ctx, id)
}

func (m *mockUserStore) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.getByEmailFn(ctx, email)
}

func (m *mockUserStore) ActivateByInvitationToken(ctx context.Context, token string) error {
	return m.activateByInvitationTokenFn(ctx, token)
}

func (m *mockUserStore) Update(ctx context.Context, p payload.UpdateUserPayload) error {
	return m.updateFn(ctx, p)
}

func (m *mockUserStore) Delete(ctx context.Context, id string) error {
	return m.deleteFn(ctx, id)
}

func (m *mockUserStore) DeleteToken(ctx context.Context, token string) error {
	return m.deleteTokenFn(ctx, token)
}

func newTestApp(users *mockUserStore) *application {
	return &application{
		config: config{mail: mailConfig{expiry: time.Hour}},
		store:  store.Storage{Users: users},
	}
}

func decodeErrorBody(t *testing.T, body *bytes.Buffer) string {
	t.Helper()
	var envelope struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body.Bytes(), &envelope); err != nil {
		t.Fatalf("failed to decode error envelope: %v (body=%s)", err, body.String())
	}
	return envelope.Error
}

func TestCreateUserHandler(t *testing.T) {
	validBody := `{"email":"jane@example.com","first_name":"Jane","last_name":"Doe","username":"jane","password":"hunter2!!"}`

	t.Run("success returns 201", func(t *testing.T) {
		users := &mockUserStore{
			createWithInvitationFn: func(ctx context.Context, u *model.User, expiry time.Duration) (*string, error) {
				if u.Email != "jane@example.com" || u.Username != "jane" {
					t.Errorf("unexpected user passed to store: %+v", u)
				}
				if u.Password == "hunter2!!" {
					t.Errorf("password was not hashed before reaching the store")
				}
				token := "abc123"
				return &token, nil
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(validBody))
		rec := httptest.NewRecorder()

		app.CreateUserHandler(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected 201, got %d (body=%s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("duplicate email returns 409 with a specific message", func(t *testing.T) {
		users := &mockUserStore{
			createWithInvitationFn: func(ctx context.Context, u *model.User, expiry time.Duration) (*string, error) {
				return nil, &pq.Error{Constraint: "users_email_key"}
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(validBody))
		rec := httptest.NewRecorder()

		app.CreateUserHandler(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d (body=%s)", rec.Code, rec.Body.String())
		}
		if got := decodeErrorBody(t, rec.Body); got != "User with this email already exists" {
			t.Errorf("unexpected error message: %q", got)
		}
	})

	t.Run("duplicate username returns 409 with a specific message", func(t *testing.T) {
		users := &mockUserStore{
			createWithInvitationFn: func(ctx context.Context, u *model.User, expiry time.Duration) (*string, error) {
				return nil, &pq.Error{Constraint: "users_username_key"}
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(validBody))
		rec := httptest.NewRecorder()

		app.CreateUserHandler(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf("expected 409, got %d (body=%s)", rec.Code, rec.Body.String())
		}
		if got := decodeErrorBody(t, rec.Body); got != "User with this username already exists" {
			t.Errorf("unexpected error message: %q", got)
		}
	})

	t.Run("unmapped constraint violation returns 500", func(t *testing.T) {
		users := &mockUserStore{
			createWithInvitationFn: func(ctx context.Context, u *model.User, expiry time.Duration) (*string, error) {
				return nil, &pq.Error{Constraint: "some_other_constraint"}
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(validBody))
		rec := httptest.NewRecorder()

		app.CreateUserHandler(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("generic store error returns 500", func(t *testing.T) {
		users := &mockUserStore{
			createWithInvitationFn: func(ctx context.Context, u *model.User, expiry time.Duration) (*string, error) {
				return nil, context.DeadlineExceeded
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(validBody))
		rec := httptest.NewRecorder()

		app.CreateUserHandler(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("malformed json returns 400", func(t *testing.T) {
		users := &mockUserStore{}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/register", bytes.NewBufferString(`{"email":`))
		rec := httptest.NewRecorder()

		app.CreateUserHandler(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d (body=%s)", rec.Code, rec.Body.String())
		}
	})
}

func hashPassword(t *testing.T, plain string) string {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	return string(hash)
}

func TestLoginHandler(t *testing.T) {
	t.Run("valid credentials for an active user return 200", func(t *testing.T) {
		users := &mockUserStore{
			getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return &model.User{ID: "u1", Email: email, Password: hashPassword(t, "correct-pw"), IsActive: true}, nil
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBufferString(`{"email":"jane@example.com","password":"correct-pw"}`))
		rec := httptest.NewRecorder()

		app.LoginHandler(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d (body=%s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("unknown email returns 401", func(t *testing.T) {
		users := &mockUserStore{
			getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return nil, store.ErrorNotFound
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBufferString(`{"email":"ghost@example.com","password":"whatever"}`))
		rec := httptest.NewRecorder()

		app.LoginHandler(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d (body=%s)", rec.Code, rec.Body.String())
		}
		if got := decodeErrorBody(t, rec.Body); got != "Invalid email or password" {
			t.Errorf("unexpected error message: %q", got)
		}
	})

	t.Run("wrong password returns 401", func(t *testing.T) {
		users := &mockUserStore{
			getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return &model.User{ID: "u1", Email: email, Password: hashPassword(t, "correct-pw"), IsActive: true}, nil
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBufferString(`{"email":"jane@example.com","password":"wrong-pw"}`))
		rec := httptest.NewRecorder()

		app.LoginHandler(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d (body=%s)", rec.Code, rec.Body.String())
		}
		if got := decodeErrorBody(t, rec.Body); got != "Invalid email or password" {
			t.Errorf("unexpected error message: %q", got)
		}
	})

	t.Run("inactive account returns 401 and does not fall through to success", func(t *testing.T) {
		users := &mockUserStore{
			getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return &model.User{ID: "u1", Email: email, Password: hashPassword(t, "correct-pw"), IsActive: false}, nil
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBufferString(`{"email":"jane@example.com","password":"correct-pw"}`))
		rec := httptest.NewRecorder()

		app.LoginHandler(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d (body=%s)", rec.Code, rec.Body.String())
		}
		if got := decodeErrorBody(t, rec.Body); got != "Not an active user" {
			t.Errorf("unexpected error message: %q", got)
		}
		// A pre-existing bug let this handler write a second (200) response
		// after the 401. Guard against a regression: exactly one write.
		if rec.Result().ContentLength == 0 && rec.Body.Len() == 0 {
			t.Fatalf("expected a response body")
		}
	})

	t.Run("store error returns 500", func(t *testing.T) {
		users := &mockUserStore{
			getByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
				return nil, context.DeadlineExceeded
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/login", bytes.NewBufferString(`{"email":"jane@example.com","password":"whatever"}`))
		rec := httptest.NewRecorder()

		app.LoginHandler(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%s)", rec.Code, rec.Body.String())
		}
	})
}

func TestActivateUserHandler(t *testing.T) {
	t.Run("valid token returns 202", func(t *testing.T) {
		users := &mockUserStore{
			activateByInvitationTokenFn: func(ctx context.Context, token string) error {
				if token != "tok-123" {
					t.Errorf("unexpected token: %q", token)
				}
				return nil
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/activate/tok-123", nil)
		rec := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("token", "tok-123")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		app.ActivateUserHandler(rec, req)

		if rec.Code != http.StatusAccepted {
			t.Fatalf("expected 202, got %d (body=%s)", rec.Code, rec.Body.String())
		}
	})

	t.Run("store error returns 500 and does not also write success", func(t *testing.T) {
		users := &mockUserStore{
			activateByInvitationTokenFn: func(ctx context.Context, token string) error {
				return store.ErrorNotFound
			},
		}
		app := newTestApp(users)

		req := httptest.NewRequest(http.MethodPost, "/v1/auth/activate/bad-token", nil)
		rec := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("token", "bad-token")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		app.ActivateUserHandler(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d (body=%s)", rec.Code, rec.Body.String())
		}
	})
}

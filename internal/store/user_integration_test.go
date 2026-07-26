//go:build integration

package store

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/eedriz99/go_blog/internal/dto/payload"
	"github.com/eedriz99/go_blog/internal/model"
	"github.com/eedriz99/go_blog/internal/testutil"
	"github.com/lib/pq"
)

var userSeq atomic.Int64

// newUniqueUser builds a user with a collision-free email/username so
// tests can run repeatedly against a shared, already-seeded dev database
// without clobbering existing rows.
func newUniqueUser(t *testing.T) *model.User {
	t.Helper()
	suffix := fmt.Sprintf("%d-%d", time.Now().UnixNano(), userSeq.Add(1))
	return &model.User{
		Email:     fmt.Sprintf("itest-%s@example.com", suffix),
		FirstName: "Integration",
		LastName:  "Tester",
		Username:  fmt.Sprintf("itest-%s", suffix),
		Password:  "bcrypt-hash-placeholder",
	}
}

func newTestUserStore(t *testing.T) (Storage, func(email string)) {
	t.Helper()
	db := testutil.OpenTestDB(t)
	s := NewStore(db)

	cleanup := func(email string) {
		t.Cleanup(func() {
			// user_invitations.user_id has ON DELETE CASCADE, so deleting
			// the user row is enough to remove any invitation rows too.
			if _, err := db.Exec(`DELETE FROM users WHERE email = $1`, email); err != nil {
				t.Logf("cleanup: failed to delete user %s: %v", email, err)
			}
		})
	}

	return s, cleanup
}

func TestIntegration_UserStore_CreateWithInvitation(t *testing.T) {
	s, cleanup := newTestUserStore(t)
	ctx := context.Background()

	user := newUniqueUser(t)
	cleanup(user.Email)

	token, err := s.Users.CreateWithInvitation(ctx, user, time.Hour)
	if err != nil {
		t.Fatalf("CreateWithInvitation: %v", err)
	}
	if token == nil || *token == "" {
		t.Fatalf("expected a non-empty invitation token, got %v", token)
	}
	if user.ID == "" {
		t.Fatalf("expected the user's ID to be populated after creation")
	}

	stored, err := s.Users.GetByEmail(ctx, user.Email)
	if err != nil {
		t.Fatalf("GetByEmail: %v", err)
	}
	if stored.Username != user.Username {
		t.Errorf("username = %q, want %q", stored.Username, user.Username)
	}
	if stored.Password != user.Password {
		t.Errorf("password hash was not stored as-is")
	}
}

func TestIntegration_UserStore_CreateWithInvitation_DuplicateEmail(t *testing.T) {
	s, cleanup := newTestUserStore(t)
	ctx := context.Background()

	user := newUniqueUser(t)
	cleanup(user.Email)

	if _, err := s.Users.CreateWithInvitation(ctx, user, time.Hour); err != nil {
		t.Fatalf("first CreateWithInvitation: %v", err)
	}

	dupe := newUniqueUser(t)
	dupe.Email = user.Email // same email, different username

	_, err := s.Users.CreateWithInvitation(ctx, dupe, time.Hour)
	if err == nil {
		t.Fatalf("expected a unique-violation error, got nil")
	}

	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		t.Fatalf("expected a *pq.Error, got %T: %v", err, err)
	}
	if pqErr.Constraint != "users_email_key" {
		t.Errorf("constraint = %q, want %q", pqErr.Constraint, "users_email_key")
	}
}

func TestIntegration_UserStore_GetByEmail_NotFound(t *testing.T) {
	s, _ := newTestUserStore(t)
	ctx := context.Background()

	_, err := s.Users.GetByEmail(ctx, "does-not-exist-xyz@example.com")
	if !errors.Is(err, ErrorNotFound) {
		t.Fatalf("expected ErrorNotFound, got %v", err)
	}
}

func TestIntegration_UserStore_ActivateByInvitationToken(t *testing.T) {
	s, cleanup := newTestUserStore(t)
	ctx := context.Background()

	user := newUniqueUser(t)
	cleanup(user.Email)

	token, err := s.Users.CreateWithInvitation(ctx, user, time.Hour)
	if err != nil {
		t.Fatalf("CreateWithInvitation: %v", err)
	}

	before, err := s.Users.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID before activation: %v", err)
	}
	if before.IsActive {
		t.Fatalf("expected a freshly created user to be inactive")
	}

	if err := s.Users.ActivateByInvitationToken(ctx, *token); err != nil {
		t.Fatalf("ActivateByInvitationToken: %v", err)
	}

	activated, err := s.Users.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID after activation: %v", err)
	}
	if activated.ID != user.ID {
		t.Fatalf("fetched wrong user")
	}
	if !activated.IsActive {
		t.Fatalf("expected user to be active after ActivateByInvitationToken, got IsActive=false")
	}

	// The token should be single-use: activating again (or deleting it
	// directly) must now report it as gone.
	if err := s.Users.DeleteToken(ctx, *token); !errors.Is(err, ErrorNotFound) {
		t.Errorf("expected the token to already be consumed, got err=%v", err)
	}
}

func TestIntegration_UserStore_ActivateByInvitationToken_InvalidToken(t *testing.T) {
	s, _ := newTestUserStore(t)
	ctx := context.Background()

	err := s.Users.ActivateByInvitationToken(ctx, "00000000-0000-0000-0000-000000000000")
	if !errors.Is(err, ErrorNotFound) {
		t.Fatalf("expected ErrorNotFound for an unknown token, got %v", err)
	}
}

func TestIntegration_UserStore_ActivateByInvitationToken_ExpiredToken(t *testing.T) {
	s, cleanup := newTestUserStore(t)
	ctx := context.Background()

	user := newUniqueUser(t)
	cleanup(user.Email)

	// A negative expiry places the invitation's expiry timestamp in the
	// past, exercising the same "already expired" path a real token would
	// hit after its window elapses.
	token, err := s.Users.CreateWithInvitation(ctx, user, -time.Hour)
	if err != nil {
		t.Fatalf("CreateWithInvitation: %v", err)
	}

	err = s.Users.ActivateByInvitationToken(ctx, *token)
	if !errors.Is(err, ErrorNotFound) {
		t.Fatalf("expected ErrorNotFound for an expired token, got %v", err)
	}
}

func TestIntegration_UserStore_Update(t *testing.T) {
	s, cleanup := newTestUserStore(t)
	ctx := context.Background()

	user := newUniqueUser(t)
	cleanup(user.Email)

	if _, err := s.Users.CreateWithInvitation(ctx, user, time.Hour); err != nil {
		t.Fatalf("CreateWithInvitation: %v", err)
	}

	newFirstName := "Updated"
	err := s.Users.Update(ctx, payload.UpdateUserPayload{
		ID:        user.ID,
		FirstName: &newFirstName,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := s.Users.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.FirstName != newFirstName {
		t.Errorf("first_name = %q, want %q", got.FirstName, newFirstName)
	}
	if got.LastName != user.LastName {
		t.Errorf("last_name changed unexpectedly: got %q, want unchanged %q", got.LastName, user.LastName)
	}
}

func TestIntegration_UserStore_Delete(t *testing.T) {
	s, cleanup := newTestUserStore(t)
	ctx := context.Background()

	user := newUniqueUser(t)
	cleanup(user.Email)

	if _, err := s.Users.CreateWithInvitation(ctx, user, time.Hour); err != nil {
		t.Fatalf("CreateWithInvitation: %v", err)
	}

	if err := s.Users.Delete(ctx, user.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if _, err := s.Users.GetByID(ctx, user.ID); !errors.Is(err, ErrorNotFound) {
		t.Fatalf("expected ErrorNotFound after delete, got %v", err)
	}

	if err := s.Users.Delete(ctx, user.ID); !errors.Is(err, ErrorNotFound) {
		t.Fatalf("expected ErrorNotFound deleting an already-deleted user, got %v", err)
	}
}

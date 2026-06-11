package auth_test

import (
	"testing"

	"github.com/adeelkhan/qr-service/internal/auth"
	"github.com/adeelkhan/qr-service/internal/database"
)

func newTestService(t *testing.T) *auth.Service {
	t.Helper()
	db, err := database.ConnectSQLite()
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	return auth.NewService(db, "test-secret")
}

func TestRegister_HappyPath(t *testing.T) {
	svc := newTestService(t)
	user, err := svc.Register("alice", "alice@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != "alice" {
		t.Errorf("got username %q, want %q", user.Username, "alice")
	}
	if user.PasswordHash == "password123" {
		t.Error("password must be hashed")
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.Register("bob", "bob@example.com", "pass"); err != nil {
		t.Fatalf("first register: %v", err)
	}
	_, err := svc.Register("bob", "bob2@example.com", "pass")
	if err == nil {
		t.Fatal("expected error for duplicate username, got nil")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.Register("carol", "carol@example.com", "pass"); err != nil {
		t.Fatalf("first register: %v", err)
	}
	_, err := svc.Register("carol2", "carol@example.com", "pass")
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
}

func TestLogin_ValidCredentials(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.Register("dave", "dave@example.com", "secret"); err != nil {
		t.Fatalf("register: %v", err)
	}
	token, expiresAt, err := svc.Login("dave", "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	if expiresAt.IsZero() {
		t.Error("expected non-zero expiresAt")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := newTestService(t)
	if _, err := svc.Register("eve", "eve@example.com", "correct"); err != nil {
		t.Fatalf("register: %v", err)
	}
	_, _, err := svc.Login("eve", "wrong")
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	svc := newTestService(t)
	_, _, err := svc.Login("nobody", "pass")
	if err == nil {
		t.Fatal("expected error for unknown user, got nil")
	}
}

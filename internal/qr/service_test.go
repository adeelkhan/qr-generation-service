package qr_test

import (
	"testing"

	"github.com/adeelkhan/qr-service/internal/database"
	"github.com/adeelkhan/qr-service/internal/models"
	"github.com/adeelkhan/qr-service/internal/qr"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setup(t *testing.T) (*qr.Service, *gorm.DB) {
	t.Helper()
	db, err := database.ConnectSQLite()
	if err != nil {
		t.Fatalf("sqlite: %v", err)
	}
	return qr.NewService(db), db
}

func seedUser(t *testing.T, db *gorm.DB) models.User {
	t.Helper()
	u := models.User{Username: "tester", Email: "t@t.com", PasswordHash: "x"}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return u
}

func TestGenerate_CreatesRecord(t *testing.T) {
	svc, db := setup(t)
	user := seedUser(t, db)
	code, err := svc.Generate(user.ID, "hello world", "My QR")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(code.ImageData) == 0 {
		t.Error("expected non-empty image data")
	}
	if code.InputText != "hello world" {
		t.Errorf("got input %q, want %q", code.InputText, "hello world")
	}
}

func TestList_ReturnsOnlyOwnerCodes(t *testing.T) {
	svc, db := setup(t)
	u1 := seedUser(t, db)
	u2 := models.User{Username: "other", Email: "o@o.com", PasswordHash: "x"}
	db.Create(&u2)

	svc.Generate(u1.ID, "for u1", "")
	svc.Generate(u2.ID, "for u2", "")

	codes, err := svc.List(u1.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(codes) != 1 {
		t.Errorf("got %d codes, want 1", len(codes))
	}
	if codes[0].UserID != u1.ID {
		t.Error("returned code belongs to wrong user")
	}
}

func TestGet_HappyPath(t *testing.T) {
	svc, db := setup(t)
	user := seedUser(t, db)
	created, _ := svc.Generate(user.ID, "test", "")
	got, err := svc.Get(user.ID, created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Error("returned wrong record")
	}
}

func TestGet_WrongUser(t *testing.T) {
	svc, db := setup(t)
	owner := seedUser(t, db)
	other := models.User{Username: "intruder", Email: "i@i.com", PasswordHash: "x"}
	db.Create(&other)
	created, _ := svc.Generate(owner.ID, "secret", "")
	_, err := svc.Get(other.ID, created.ID)
	if err == nil {
		t.Fatal("expected error when wrong user fetches QR code")
	}
}

func TestGet_NotFound(t *testing.T) {
	svc, db := setup(t)
	user := seedUser(t, db)
	_, err := svc.Get(user.ID, uuid.New())
	if err == nil {
		t.Fatal("expected error for non-existent QR code")
	}
}

func TestDelete_HappyPath(t *testing.T) {
	svc, db := setup(t)
	user := seedUser(t, db)
	created, _ := svc.Generate(user.ID, "bye", "")
	if err := svc.Delete(user.ID, created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := svc.Get(user.ID, created.ID)
	if err == nil {
		t.Fatal("expected record to be gone after delete")
	}
}

func TestDelete_WrongUser(t *testing.T) {
	svc, db := setup(t)
	owner := seedUser(t, db)
	other := models.User{Username: "attacker", Email: "a@a.com", PasswordHash: "x"}
	db.Create(&other)
	created, _ := svc.Generate(owner.ID, "mine", "")
	err := svc.Delete(other.ID, created.ID)
	if err == nil {
		t.Fatal("expected error when wrong user deletes QR code")
	}
}

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/isa0-gh/reader/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
	users  map[uint]*model.User
	nextID uint
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{users: map[uint]*model.User{}}
}

func (f *fakeUserRepo) Create(_ context.Context, u *model.User) error {
	f.nextID++
	u.ID = f.nextID
	cp := *u
	f.users[u.ID] = &cp
	return nil
}

func (f *fakeUserRepo) GetByID(_ context.Context, id uint) (*model.User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	cp := *u
	return &cp, nil
}

func (f *fakeUserRepo) GetByEmail(_ context.Context, email string) (*model.User, error) {
	for _, u := range f.users {
		if u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, errors.New("user not found")
}

func (f *fakeUserRepo) Update(_ context.Context, u *model.User) error {
	if _, ok := f.users[u.ID]; !ok {
		return errors.New("user not found")
	}
	cp := *u
	f.users[u.ID] = &cp
	return nil
}

func (f *fakeUserRepo) Delete(_ context.Context, id uint) error {
	if _, ok := f.users[id]; !ok {
		return errors.New("user not found")
	}
	delete(f.users, id)
	return nil
}

func (f *fakeUserRepo) List(_ context.Context, limit int, after, before uint) ([]model.User, error) {
	var list []model.User
	for _, u := range f.users {
		list = append(list, *u)
	}
	return list, nil
}

func (f *fakeUserRepo) Count(_ context.Context) (int64, error) {
	return int64(len(f.users)), nil
}

func (f *fakeUserRepo) SetAvatar(_ context.Context, userID uint, avatar *model.S3Object) (*model.User, error) {
	u, ok := f.users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	avatar.ID = 1
	u.Avatar = avatar
	cp := *u
	return &cp, nil
}

func newTestUser(repo *fakeUserRepo, password string) *model.User {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &model.User{Email: "a@b.com", PasswordHash: string(hash), JwtID: "original-jti"}
	repo.Create(context.Background(), u)
	return u
}

func TestChangePassword_WrongCurrentPassword(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)
	u := newTestUser(repo, "correct-horse")

	err := svc.ChangePassword(context.Background(), u.ID, "wrong-password", "new-password")
	if err == nil {
		t.Fatal("expected error for wrong current password, got nil")
	}

	stored, _ := repo.GetByID(context.Background(), u.ID)
	if bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte("correct-horse")) != nil {
		t.Error("password hash should be unchanged after a failed attempt")
	}
}

func TestChangePassword_Success(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)
	u := newTestUser(repo, "correct-horse")

	if err := svc.ChangePassword(context.Background(), u.ID, "correct-horse", "new-password"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stored, _ := repo.GetByID(context.Background(), u.ID)
	if bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte("new-password")) != nil {
		t.Error("expected new password to verify against stored hash")
	}
	if stored.JwtID == "original-jti" {
		t.Error("expected JwtID to rotate on password change, invalidating existing sessions")
	}
}

func TestChangeEmail_WrongCurrentPassword(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)
	u := newTestUser(repo, "correct-horse")

	if _, err := svc.ChangeEmail(context.Background(), u.ID, "wrong-password", "new@b.com"); err == nil {
		t.Fatal("expected error for wrong current password, got nil")
	}

	stored, _ := repo.GetByID(context.Background(), u.ID)
	if stored.Email != "a@b.com" {
		t.Error("email should be unchanged after a failed attempt")
	}
}

func TestChangeEmail_AlreadyInUse(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)
	u := newTestUser(repo, "correct-horse")

	other := &model.User{Email: "taken@b.com", PasswordHash: u.PasswordHash}
	repo.Create(context.Background(), other)

	if _, err := svc.ChangeEmail(context.Background(), u.ID, "correct-horse", "taken@b.com"); err == nil {
		t.Fatal("expected error for email already in use, got nil")
	}
}

func TestChangeEmail_Success(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)
	u := newTestUser(repo, "correct-horse")

	updated, err := svc.ChangeEmail(context.Background(), u.ID, "correct-horse", "new@b.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Email != "new@b.com" {
		t.Errorf("got email %q, want new@b.com", updated.Email)
	}
	if updated.JwtID == "original-jti" {
		t.Error("expected JwtID to rotate on email change, invalidating existing sessions")
	}
}

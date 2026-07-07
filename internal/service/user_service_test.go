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

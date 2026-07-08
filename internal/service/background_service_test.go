package service

import (
	"context"
	"errors"
	"testing"

	"github.com/isa0-gh/reader/internal/model"
)

type fakeBackgroundRepo struct {
	backgrounds map[uint]*model.Background
	nextID      uint
}

func newFakeBackgroundRepo() *fakeBackgroundRepo {
	return &fakeBackgroundRepo{backgrounds: map[uint]*model.Background{}}
}

func (f *fakeBackgroundRepo) List(_ context.Context) ([]model.Background, error) {
	list := []model.Background{}
	for _, b := range f.backgrounds {
		list = append(list, *b)
	}
	return list, nil
}

func (f *fakeBackgroundRepo) Create(_ context.Context, b *model.Background) error {
	f.nextID++
	b.ID = f.nextID
	b.Image.ID = f.nextID
	b.ImageID = b.Image.ID
	cp := *b
	f.backgrounds[b.ID] = &cp
	return nil
}

func (f *fakeBackgroundRepo) Delete(_ context.Context, id uint) error {
	if _, ok := f.backgrounds[id]; !ok {
		return errors.New("background not found")
	}
	delete(f.backgrounds, id)
	return nil
}

func TestAddBackground(t *testing.T) {
	repo := newFakeBackgroundRepo()
	svc := NewBackgroundService(repo)

	bg, err := svc.AddBackground(context.Background(), "backgrounds/a.jpg", "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bg.ID == 0 || bg.ImageID == 0 {
		t.Errorf("expected assigned IDs, got %+v", bg)
	}
	if bg.Image.Key != "backgrounds/a.jpg" {
		t.Errorf("got key %q, want backgrounds/a.jpg", bg.Image.Key)
	}
}

func TestDeleteBackground(t *testing.T) {
	repo := newFakeBackgroundRepo()
	svc := NewBackgroundService(repo)

	bg, _ := svc.AddBackground(context.Background(), "backgrounds/a.jpg", "my-bucket")
	if err := svc.DeleteBackground(context.Background(), bg.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, err := svc.ListBackgrounds(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected no backgrounds after delete, got %d", len(list))
	}
}

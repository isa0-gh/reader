package service

import (
	"context"
	"testing"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type fakeFavoriteRepo struct {
	favorites map[uint]*model.Favorite
	nextID    uint
}

func newFakeFavoriteRepo() *fakeFavoriteRepo {
	return &fakeFavoriteRepo{favorites: map[uint]*model.Favorite{}}
}

func (f *fakeFavoriteRepo) GetByUserAndSeries(_ context.Context, userID, seriesID uint) (*model.Favorite, error) {
	for _, fav := range f.favorites {
		if fav.UserID == userID && fav.SeriesID == seriesID {
			cp := *fav
			return &cp, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeFavoriteRepo) Create(_ context.Context, fav *model.Favorite) error {
	f.nextID++
	fav.ID = f.nextID
	cp := *fav
	f.favorites[fav.ID] = &cp
	return nil
}

func (f *fakeFavoriteRepo) Delete(_ context.Context, userID, seriesID uint) error {
	for id, fav := range f.favorites {
		if fav.UserID == userID && fav.SeriesID == seriesID {
			delete(f.favorites, id)
		}
	}
	return nil
}

func (f *fakeFavoriteRepo) ListByUser(_ context.Context, userID uint, _, _ int) ([]model.Favorite, error) {
	var list []model.Favorite
	for _, fav := range f.favorites {
		if fav.UserID == userID {
			list = append(list, *fav)
		}
	}
	return list, nil
}

func (f *fakeFavoriteRepo) CountByUser(_ context.Context, userID uint) (int64, error) {
	var count int64
	for _, fav := range f.favorites {
		if fav.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (f *fakeFavoriteRepo) UpdateProgress(_ context.Context, userID, seriesID, chapterID uint) error {
	for _, fav := range f.favorites {
		if fav.UserID == userID && fav.SeriesID == seriesID {
			fav.LastReadChapterID = &chapterID
		}
	}
	return nil
}

func TestAddFavorite_Idempotent(t *testing.T) {
	repo := newFakeFavoriteRepo()
	svc := NewFavoriteService(repo)

	first, err := svc.AddFavorite(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	second, err := svc.AddFavorite(context.Background(), 1, 100)
	if err != nil {
		t.Fatalf("unexpected error on repeat favorite: %v", err)
	}

	if first.ID != second.ID {
		t.Errorf("expected favoriting twice to return the same row, got IDs %d and %d", first.ID, second.ID)
	}
	if len(repo.favorites) != 1 {
		t.Errorf("expected exactly one favorite row, got %d", len(repo.favorites))
	}
}

func TestRemoveFavorite(t *testing.T) {
	repo := newFakeFavoriteRepo()
	svc := NewFavoriteService(repo)

	svc.AddFavorite(context.Background(), 1, 100)
	if err := svc.RemoveFavorite(context.Background(), 1, 100); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, _, err := svc.ListFavorites(context.Background(), 1, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("expected no favorites after removal, got %d", len(list))
	}
}

// UpdateProgress must not error when the series isn't favorited — the
// frontend pings it on every chapter view regardless of favorite status.
func TestUpdateProgress_NoOpWhenNotFavorited(t *testing.T) {
	repo := newFakeFavoriteRepo()
	svc := NewFavoriteService(repo)

	if err := svc.UpdateProgress(context.Background(), 1, 100, 5); err != nil {
		t.Fatalf("expected no error for a non-favorited series, got %v", err)
	}
}

func TestUpdateProgress_SetsLastReadChapter(t *testing.T) {
	repo := newFakeFavoriteRepo()
	svc := NewFavoriteService(repo)

	svc.AddFavorite(context.Background(), 1, 100)
	if err := svc.UpdateProgress(context.Background(), 1, 100, 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	list, _, _ := svc.ListFavorites(context.Background(), 1, 0, 0)
	if len(list) != 1 || list[0].LastReadChapterID == nil || *list[0].LastReadChapterID != 42 {
		t.Errorf("expected last_read_chapter_id 42, got %+v", list)
	}
}

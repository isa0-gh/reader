package service

import (
	"context"
	"errors"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
	"gorm.io/gorm"
)

type FavoriteService interface {
	ListFavorites(ctx context.Context, userID uint, limit, offset int) ([]model.Favorite, int64, error)
	AddFavorite(ctx context.Context, userID, seriesID uint) (*model.Favorite, error)
	RemoveFavorite(ctx context.Context, userID, seriesID uint) error
	UpdateProgress(ctx context.Context, userID, seriesID, chapterID uint) error
}

type favoriteService struct{ repo repository.FavoriteRepository }

func NewFavoriteService(repo repository.FavoriteRepository) FavoriteService {
	return &favoriteService{repo: repo}
}

func (s *favoriteService) ListFavorites(ctx context.Context, userID uint, limit, offset int) ([]model.Favorite, int64, error) {
	list, err := s.repo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountByUser(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// AddFavorite is idempotent: favoriting an already-favorited series just
// returns the existing row instead of erroring on the unique constraint.
func (s *favoriteService) AddFavorite(ctx context.Context, userID, seriesID uint) (*model.Favorite, error) {
	existing, err := s.repo.GetByUserAndSeries(ctx, userID, seriesID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	f := &model.Favorite{UserID: userID, SeriesID: seriesID}
	if err := s.repo.Create(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *favoriteService) RemoveFavorite(ctx context.Context, userID, seriesID uint) error {
	return s.repo.Delete(ctx, userID, seriesID)
}

func (s *favoriteService) UpdateProgress(ctx context.Context, userID, seriesID, chapterID uint) error {
	return s.repo.UpdateProgress(ctx, userID, seriesID, chapterID)
}

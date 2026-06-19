package service

import (
	"context"
	"fmt"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type FavoriteService interface {
	AddFavorite(ctx context.Context, userID, seriesID uint) (*model.UserFavorite, error)
	RemoveFavorite(ctx context.Context, userID, seriesID uint) error
	GetUserFavorites(ctx context.Context, userID uint) ([]model.UserFavorite, error)
	IsFavorite(ctx context.Context, userID, seriesID uint) (bool, error)
	GetFavoriteCount(ctx context.Context, seriesID uint) (int64, error)
}

type favoriteService struct {
	repo repository.FavoriteRepository
}

func NewFavoriteService(repo repository.FavoriteRepository) FavoriteService {
	return &favoriteService{repo: repo}
}

func (s *favoriteService) AddFavorite(ctx context.Context, userID, seriesID uint) (*model.UserFavorite, error) {
	// Check if already favorited
	exists, err := s.repo.Exists(ctx, userID, seriesID)
	if err != nil {
		return nil, fmt.Errorf("check favorite: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("already favorited")
	}

	fav := &model.UserFavorite{
		UserID:   userID,
		SeriesID: seriesID,
	}
	if err := s.repo.Create(ctx, fav); err != nil {
		return nil, fmt.Errorf("create favorite: %w", err)
	}
	return fav, nil
}

func (s *favoriteService) RemoveFavorite(ctx context.Context, userID, seriesID uint) error {
	return s.repo.Delete(ctx, userID, seriesID)
}

func (s *favoriteService) GetUserFavorites(ctx context.Context, userID uint) ([]model.UserFavorite, error) {
	return s.repo.GetUserFavorites(ctx, userID)
}

func (s *favoriteService) IsFavorite(ctx context.Context, userID, seriesID uint) (bool, error) {
	return s.repo.Exists(ctx, userID, seriesID)
}

func (s *favoriteService) GetFavoriteCount(ctx context.Context, seriesID uint) (int64, error) {
	return s.repo.GetFavoriteCount(ctx, seriesID)
}

package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type FavoriteRepository interface {
	Create(ctx context.Context, fav *model.UserFavorite) error
	Delete(ctx context.Context, userID, seriesID uint) error
	GetUserFavorites(ctx context.Context, userID uint) ([]model.UserFavorite, error)
	Exists(ctx context.Context, userID, seriesID uint) (bool, error)
	GetFavoriteCount(ctx context.Context, seriesID uint) (int64, error)
}

type favoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) Create(ctx context.Context, fav *model.UserFavorite) error {
	return r.db.WithContext(ctx).Create(fav).Error
}

func (r *favoriteRepository) Delete(ctx context.Context, userID, seriesID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND series_id = ?", userID, seriesID).
		Delete(&model.UserFavorite{}).Error
}

func (r *favoriteRepository) GetUserFavorites(ctx context.Context, userID uint) ([]model.UserFavorite, error) {
	var list []model.UserFavorite
	err := r.db.WithContext(ctx).
		Preload("Series").
		Preload("Series.CoverImage").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&list).Error
	return list, err
}

func (r *favoriteRepository) Exists(ctx context.Context, userID, seriesID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.UserFavorite{}).
		Where("user_id = ? AND series_id = ?", userID, seriesID).
		Count(&count).Error
	return count > 0, err
}

func (r *favoriteRepository) GetFavoriteCount(ctx context.Context, seriesID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.UserFavorite{}).
		Where("series_id = ?", seriesID).
		Count(&count).Error
	return count, err
}

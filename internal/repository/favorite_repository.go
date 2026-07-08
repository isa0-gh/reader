package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type FavoriteRepository interface {
	GetByUserAndSeries(ctx context.Context, userID, seriesID uint) (*model.Favorite, error)
	Create(ctx context.Context, f *model.Favorite) error
	Delete(ctx context.Context, userID, seriesID uint) error
	ListByUser(ctx context.Context, userID uint, limit, offset int) ([]model.Favorite, error)
	CountByUser(ctx context.Context, userID uint) (int64, error)
	UpdateProgress(ctx context.Context, userID, seriesID, chapterID uint) error
}

type favoriteRepository struct{ db *gorm.DB }

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) GetByUserAndSeries(ctx context.Context, userID, seriesID uint) (*model.Favorite, error) {
	var f model.Favorite
	if err := r.db.WithContext(ctx).Where("user_id = ? AND series_id = ?", userID, seriesID).First(&f).Error; err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *favoriteRepository) Create(ctx context.Context, f *model.Favorite) error {
	return r.db.WithContext(ctx).Create(f).Error
}

func (r *favoriteRepository) Delete(ctx context.Context, userID, seriesID uint) error {
	return r.db.WithContext(ctx).Where("user_id = ? AND series_id = ?", userID, seriesID).Delete(&model.Favorite{}).Error
}

func (r *favoriteRepository) ListByUser(ctx context.Context, userID uint, limit, offset int) ([]model.Favorite, error) {
	var list []model.Favorite
	query := r.db.WithContext(ctx).
		Preload("Series").Preload("Series.CoverImage").Preload("LastReadChapter").
		Where("user_id = ?", userID).Order("updated_at desc")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *favoriteRepository) CountByUser(ctx context.Context, userID uint) (int64, error) {
	var count int64
	return count, r.db.WithContext(ctx).Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&count).Error
}

// UpdateProgress is a no-op (0 rows affected, no error) when the series isn't
// favorited — callers ping this opportunistically on every chapter view and
// shouldn't have to check favorite status first.
func (r *favoriteRepository) UpdateProgress(ctx context.Context, userID, seriesID, chapterID uint) error {
	return r.db.WithContext(ctx).Model(&model.Favorite{}).
		Where("user_id = ? AND series_id = ?", userID, seriesID).
		Update("last_read_chapter_id", chapterID).Error
}

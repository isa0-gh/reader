package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type ProgressRepository interface {
	Upsert(ctx context.Context, p *model.ReadingProgress) error
	GetByUserAndChapter(ctx context.Context, userID, chapterID uint) (*model.ReadingProgress, error)
	GetUserProgress(ctx context.Context, userID uint) ([]model.ReadingProgress, error)
	GetSeriesProgress(ctx context.Context, userID, seriesID uint) ([]model.ReadingProgress, error)
	DeleteByChapter(ctx context.Context, chapterID uint) error
}

type progressRepository struct{ db *gorm.DB }

func NewProgressRepository(db *gorm.DB) ProgressRepository {
	return &progressRepository{db: db}
}

func (r *progressRepository) Upsert(ctx context.Context, p *model.ReadingProgress) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *progressRepository) GetByUserAndChapter(ctx context.Context, userID, chapterID uint) (*model.ReadingProgress, error) {
	var p model.ReadingProgress
	if err := r.db.WithContext(ctx).Where("user_id = ? AND chapter_id = ?", userID, chapterID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *progressRepository) GetUserProgress(ctx context.Context, userID uint) ([]model.ReadingProgress, error) {
	var list []model.ReadingProgress
	err := r.db.WithContext(ctx).
		Preload("Chapter").
		Preload("Chapter.Series").
		Where("user_id = ?", userID).
		Order("updated_at desc").
		Find(&list).Error
	return list, err
}

func (r *progressRepository) GetSeriesProgress(ctx context.Context, userID, seriesID uint) ([]model.ReadingProgress, error) {
	var list []model.ReadingProgress
	err := r.db.WithContext(ctx).
		Joins("JOIN chapters ON chapters.id = reading_progresses.chapter_id").
		Where("reading_progresses.user_id = ? AND chapters.series_id = ?", userID, seriesID).
		Preload("Chapter").
		Find(&list).Error
	return list, err
}

func (r *progressRepository) DeleteByChapter(ctx context.Context, chapterID uint) error {
	return r.db.WithContext(ctx).Where("chapter_id = ?", chapterID).Delete(&model.ReadingProgress{}).Error
}

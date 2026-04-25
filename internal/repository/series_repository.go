package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type SeriesRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Series, error)
	Create(ctx context.Context, s *model.Series) error
	List(ctx context.Context) ([]model.Series, error)
	Delete(ctx context.Context, id uint) error
}

type seriesRepository struct{ db *gorm.DB }

func NewSeriesRepository(db *gorm.DB) SeriesRepository {
	return &seriesRepository{db: db}
}

func (r *seriesRepository) GetByID(ctx context.Context, id uint) (*model.Series, error) {
	var s model.Series
	if err := r.db.WithContext(ctx).Preload("Chapters").Preload("CoverImage").First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *seriesRepository) Create(ctx context.Context, s *model.Series) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *seriesRepository) List(ctx context.Context) ([]model.Series, error) {
	var list []model.Series
	if err := r.db.WithContext(ctx).Preload("CoverImage").Order("created_at desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *seriesRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Series{}, id).Error
}

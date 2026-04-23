package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type SeriesRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Series, error)
}

type seriesRepository struct{ db *gorm.DB }

func NewSeriesRepository(db *gorm.DB) SeriesRepository {
	return &seriesRepository{db: db}
}

func (r *seriesRepository) GetByID(ctx context.Context, id uint) (*model.Series, error) {
	var s model.Series
	if err := r.db.WithContext(ctx).Preload("Chapters").First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

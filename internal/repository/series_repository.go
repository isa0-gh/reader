package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type SeriesRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Series, error)
	Create(ctx context.Context, s *model.Series) error
	List(ctx context.Context, query, sort string) ([]model.Series, error)
	Delete(ctx context.Context, id uint) error
}

type seriesRepository struct{ db *gorm.DB }

func NewSeriesRepository(db *gorm.DB) SeriesRepository {
	return &seriesRepository{db: db}
}

func (r *seriesRepository) GetByID(ctx context.Context, id uint) (*model.Series, error) {
	var s model.Series
	if err := r.db.WithContext(ctx).Preload("CoverImage").First(&s, id).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Where("series_id = ?", id).Order("number asc").Find(&s.Chapters).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *seriesRepository) Create(ctx context.Context, s *model.Series) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if s.CoverImage != nil {
			if err := tx.Create(s.CoverImage).Error; err != nil {
				return err
			}
			s.CoverImageID = &s.CoverImage.ID
		}
		return tx.Omit("CoverImage").Create(s).Error
	})
}

func (r *seriesRepository) List(ctx context.Context, query, sort string) ([]model.Series, error) {
	var list []model.Series
	q := r.db.WithContext(ctx).Preload("CoverImage")
	
	if query != "" {
		q = q.Where("title ILIKE ?", "%"+query+"%")
	}
	
	switch sort {
	case "title":
		q = q.Order("title asc")
	case "created":
		q = q.Order("created_at asc")
	default:
		q = q.Order("created_at desc")
	}
	
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *seriesRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Series{}, id).Error
}

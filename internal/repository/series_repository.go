package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type SeriesRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Series, error)
	Create(ctx context.Context, s *model.Series) error
	Update(ctx context.Context, s *model.Series) error
	List(ctx context.Context, limit, offset int, q string) ([]model.Series, error)
	Count(ctx context.Context, q string) (int64, error)
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

func (r *seriesRepository) Update(ctx context.Context, s *model.Series) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cols := []string{"Title", "Slug", "Description", "Author", "Artist", "Status"}
		if s.CoverImage != nil {
			if err := tx.Create(s.CoverImage).Error; err != nil {
				return err
			}
			s.CoverImageID = &s.CoverImage.ID
			cols = append(cols, "CoverImageID")
		}
		return tx.Model(&model.Series{}).Where("id = ?", s.ID).Select(cols).Updates(s).Error
	})
}

func (r *seriesRepository) List(ctx context.Context, limit, offset int, q string) ([]model.Series, error) {
	var list []model.Series
	query := r.db.WithContext(ctx).Preload("CoverImage").Order("created_at desc")
	if q != "" {
		query = query.Where("title ILIKE ?", "%"+q+"%")
	}
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

func (r *seriesRepository) Count(ctx context.Context, q string) (int64, error) {
	query := r.db.WithContext(ctx).Model(&model.Series{})
	if q != "" {
		query = query.Where("title ILIKE ?", "%"+q+"%")
	}
	var count int64
	return count, query.Count(&count).Error
}

func (r *seriesRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Series{}, id).Error
}

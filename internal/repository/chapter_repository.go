package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type ChapterRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Chapter, error)
	Create(ctx context.Context, c *model.Chapter) error
}

type chapterRepository struct{ db *gorm.DB }

func NewChapterRepository(db *gorm.DB) ChapterRepository {
	return &chapterRepository{db: db}
}

func (r *chapterRepository) GetByID(ctx context.Context, id uint) (*model.Chapter, error) {
	var c model.Chapter
	if err := r.db.WithContext(ctx).Preload("Pages").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *chapterRepository) Create(ctx context.Context, c *model.Chapter) error {
	return r.db.WithContext(ctx).Create(c).Error
}

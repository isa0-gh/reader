package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type ChapterRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Chapter, error)
	Create(ctx context.Context, c *model.Chapter) error
	AddPages(ctx context.Context, pages []model.S3Object) error
	DeletePage(ctx context.Context, chapterID, pageID uint) error
	Delete(ctx context.Context, id uint) error
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

func (r *chapterRepository) AddPages(ctx context.Context, pages []model.S3Object) error {
	return r.db.WithContext(ctx).Create(&pages).Error
}

func (r *chapterRepository) DeletePage(ctx context.Context, chapterID, pageID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND chapter_id = ?", pageID, chapterID).Delete(&model.S3Object{}).Error
}

func (r *chapterRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Chapter{}, id).Error
}

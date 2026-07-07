package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type ChapterRepository interface {
	GetByID(ctx context.Context, id uint) (*model.Chapter, error)
	Create(ctx context.Context, c *model.Chapter) error
	Update(ctx context.Context, c *model.Chapter) error
	AddPages(ctx context.Context, pages []model.S3Object) error
	UpdatePageNumbers(ctx context.Context, chapterID uint, numbers map[uint]int) error
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

func (r *chapterRepository) Update(ctx context.Context, c *model.Chapter) error {
	return r.db.WithContext(ctx).Model(&model.Chapter{}).Where("id = ?", c.ID).
		Select("Number", "Title").Updates(c).Error
}

func (r *chapterRepository) AddPages(ctx context.Context, pages []model.S3Object) error {
	return r.db.WithContext(ctx).Create(&pages).Error
}

func (r *chapterRepository) UpdatePageNumbers(ctx context.Context, chapterID uint, numbers map[uint]int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for pageID, num := range numbers {
			if err := tx.Model(&model.S3Object{}).
				Where("id = ? AND chapter_id = ?", pageID, chapterID).
				Update("page_number", num).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *chapterRepository) DeletePage(ctx context.Context, chapterID, pageID uint) error {
	return r.db.WithContext(ctx).Where("id = ? AND chapter_id = ?", pageID, chapterID).Delete(&model.S3Object{}).Error
}

func (r *chapterRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Chapter{}, id).Error
}

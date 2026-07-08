package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(ctx context.Context, c *model.Comment) error
	GetByID(ctx context.Context, id uint) (*model.Comment, error)
	ListByChapter(ctx context.Context, chapterID uint, limit, offset int) ([]model.Comment, error)
	CountByChapter(ctx context.Context, chapterID uint) (int64, error)
	ListByUser(ctx context.Context, userID uint, limit, offset int) ([]model.Comment, error)
	CountByUser(ctx context.Context, userID uint) (int64, error)
	Delete(ctx context.Context, id uint) error
}

type commentRepository struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, c *model.Comment) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *commentRepository) GetByID(ctx context.Context, id uint) (*model.Comment, error) {
	var c model.Comment
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *commentRepository) ListByChapter(ctx context.Context, chapterID uint, limit, offset int) ([]model.Comment, error) {
	// Initialized (not nil) so json.Marshal always encodes "[]", never
	// "null" — the frontend reads res.items unconditionally.
	list := []model.Comment{}
	query := r.db.WithContext(ctx).Preload("User").Preload("User.Avatar").Where("chapter_id = ?", chapterID).Order("created_at asc")
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

func (r *commentRepository) CountByChapter(ctx context.Context, chapterID uint) (int64, error) {
	var count int64
	return count, r.db.WithContext(ctx).Model(&model.Comment{}).Where("chapter_id = ?", chapterID).Count(&count).Error
}

func (r *commentRepository) ListByUser(ctx context.Context, userID uint, limit, offset int) ([]model.Comment, error) {
	// Initialized (not nil) so json.Marshal always encodes "[]", never
	// "null" — the frontend reads res.items unconditionally.
	list := []model.Comment{}
	query := r.db.WithContext(ctx).Preload("Chapter").Where("user_id = ?", userID).Order("created_at desc")
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

func (r *commentRepository) CountByUser(ctx context.Context, userID uint) (int64, error) {
	var count int64
	return count, r.db.WithContext(ctx).Model(&model.Comment{}).Where("user_id = ?", userID).Count(&count).Error
}

func (r *commentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id uint) error
	SoftDelete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*model.Comment, error)
	GetByChapter(ctx context.Context, chapterID uint, limit, offset int) ([]model.Comment, int64, error)
	GetReplies(ctx context.Context, parentID uint, limit, offset int) ([]model.Comment, int64, error)
	IncrementReplyCount(ctx context.Context, parentID uint) error
	DecrementReplyCount(ctx context.Context, parentID uint) error
	
	// Likes
	AddLike(ctx context.Context, like *model.CommentLike) error
	RemoveLike(ctx context.Context, commentID, userID uint) error
	HasLiked(ctx context.Context, commentID, userID uint) (bool, error)
	IncrementLikeCount(ctx context.Context, commentID uint) error
	DecrementLikeCount(ctx context.Context, commentID uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) Update(ctx context.Context, comment *model.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Comment{}, id).Error
}

func (r *commentRepository) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted":    true,
			"content":    "[deleted]",
			"deleted_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *commentRepository) GetByID(ctx context.Context, id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("id = ? AND deleted = ?", id, false).
		First(&comment).Error
	return &comment, err
}

func (r *commentRepository) GetByChapter(ctx context.Context, chapterID uint, limit, offset int) ([]model.Comment, int64, error) {
	var comments []model.Comment
	var total int64
	
	query := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("chapter_id = ? AND parent_id IS NULL AND deleted = ?", chapterID, false)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := query.
		Preload("User").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&comments).Error
	
	return comments, total, err
}

func (r *commentRepository) GetReplies(ctx context.Context, parentID uint, limit, offset int) ([]model.Comment, int64, error) {
	var replies []model.Comment
	var total int64
	
	query := r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("parent_id = ? AND deleted = ?", parentID, false)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := query.
		Preload("User").
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&replies).Error
	
	return replies, total, err
}

func (r *commentRepository) IncrementReplyCount(ctx context.Context, parentID uint) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ?", parentID).
		Update("reply_count", gorm.Expr("reply_count + 1")).Error
}

func (r *commentRepository) DecrementReplyCount(ctx context.Context, parentID uint) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ? AND reply_count > 0", parentID).
		Update("reply_count", gorm.Expr("reply_count - 1")).Error
}

func (r *commentRepository) AddLike(ctx context.Context, like *model.CommentLike) error {
	return r.db.WithContext(ctx).Create(like).Error
}

func (r *commentRepository) RemoveLike(ctx context.Context, commentID, userID uint) error {
	return r.db.WithContext(ctx).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Delete(&model.CommentLike{}).Error
}

func (r *commentRepository) HasLiked(ctx context.Context, commentID, userID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.CommentLike{}).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		Count(&count).Error
	return count > 0, err
}

func (r *commentRepository) IncrementLikeCount(ctx context.Context, commentID uint) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ?", commentID).
		Update("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *commentRepository) DecrementLikeCount(ctx context.Context, commentID uint) error {
	return r.db.WithContext(ctx).Model(&model.Comment{}).
		Where("id = ? AND like_count > 0", commentID).
		Update("like_count", gorm.Expr("like_count - 1")).Error
}

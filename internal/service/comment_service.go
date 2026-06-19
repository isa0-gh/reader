package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type CommentService interface {
	CreateComment(ctx context.Context, userID, chapterID uint, content string, parentID *uint) (*model.Comment, error)
	UpdateComment(ctx context.Context, commentID, userID uint, content string) (*model.Comment, error)
	DeleteComment(ctx context.Context, commentID, userID uint, isAdmin bool) error
	GetComment(ctx context.Context, commentID uint) (*model.Comment, error)
	GetChapterComments(ctx context.Context, chapterID uint, page, pageSize int) ([]model.Comment, int64, error)
	GetCommentReplies(ctx context.Context, parentID uint, page, pageSize int) ([]model.Comment, int64, error)
	
	LikeComment(ctx context.Context, commentID, userID uint) error
	UnlikeComment(ctx context.Context, commentID, userID uint) error
	HasLiked(ctx context.Context, commentID, userID uint) (bool, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (s *commentService) CreateComment(ctx context.Context, userID, chapterID uint, content string, parentID *uint) (*model.Comment, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}
	if len(content) > 2000 {
		return nil, fmt.Errorf("comment too long (max 2000 characters)")
	}
	
	comment := &model.Comment{
		UserID:    userID,
		ChapterID: chapterID,
		Content:   content,
		ParentID:  parentID,
	}
	
	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	
	// Increment parent's reply count if this is a reply
	if parentID != nil {
		if err := s.repo.IncrementReplyCount(ctx, *parentID); err != nil {
			// Log but don't fail
			fmt.Printf("failed to increment reply count: %v\n", err)
		}
	}
	
	return comment, nil
}

func (s *commentService) UpdateComment(ctx context.Context, commentID, userID uint, content string) (*model.Comment, error) {
	comment, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		return nil, fmt.Errorf("comment not found")
	}
	
	if comment.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}
	
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}
	if len(content) > 2000 {
		return nil, fmt.Errorf("comment too long (max 2000 characters)")
	}
	
	comment.Content = content
	comment.Edited = true
	
	if err := s.repo.Update(ctx, comment); err != nil {
		return nil, fmt.Errorf("update comment: %w", err)
	}
	
	return comment, nil
}

func (s *commentService) DeleteComment(ctx context.Context, commentID, userID uint, isAdmin bool) error {
	comment, err := s.repo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("comment not found")
	}
	
	if comment.UserID != userID && !isAdmin {
		return fmt.Errorf("unauthorized")
	}
	
	// Soft delete if there are replies, hard delete otherwise
	if comment.ReplyCount > 0 {
		if err := s.repo.SoftDelete(ctx, commentID); err != nil {
			return fmt.Errorf("delete comment: %w", err)
		}
	} else {
		if err := s.repo.Delete(ctx, commentID); err != nil {
			return fmt.Errorf("delete comment: %w", err)
		}
	}
	
	// Decrement parent's reply count if this was a reply
	if comment.ParentID != nil {
		if err := s.repo.DecrementReplyCount(ctx, *comment.ParentID); err != nil {
			fmt.Printf("failed to decrement reply count: %v\n", err)
		}
	}
	
	return nil
}

func (s *commentService) GetComment(ctx context.Context, commentID uint) (*model.Comment, error) {
	return s.repo.GetByID(ctx, commentID)
}

func (s *commentService) GetChapterComments(ctx context.Context, chapterID uint, page, pageSize int) ([]model.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	offset := (page - 1) * pageSize
	return s.repo.GetByChapter(ctx, chapterID, pageSize, offset)
}

func (s *commentService) GetCommentReplies(ctx context.Context, parentID uint, page, pageSize int) ([]model.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	
	offset := (page - 1) * pageSize
	return s.repo.GetReplies(ctx, parentID, pageSize, offset)
}

func (s *commentService) LikeComment(ctx context.Context, commentID, userID uint) error {
	// Check if already liked
	liked, err := s.repo.HasLiked(ctx, commentID, userID)
	if err != nil {
		return fmt.Errorf("check like: %w", err)
	}
	if liked {
		return fmt.Errorf("already liked")
	}
	
	like := &model.CommentLike{
		CommentID: commentID,
		UserID:    userID,
	}
	
	if err := s.repo.AddLike(ctx, like); err != nil {
		return fmt.Errorf("add like: %w", err)
	}
	
	if err := s.repo.IncrementLikeCount(ctx, commentID); err != nil {
		fmt.Printf("failed to increment like count: %v\n", err)
	}
	
	return nil
}

func (s *commentService) UnlikeComment(ctx context.Context, commentID, userID uint) error {
	if err := s.repo.RemoveLike(ctx, commentID, userID); err != nil {
		return fmt.Errorf("remove like: %w", err)
	}
	
	if err := s.repo.DecrementLikeCount(ctx, commentID); err != nil {
		fmt.Printf("failed to decrement like count: %v\n", err)
	}
	
	return nil
}

func (s *commentService) HasLiked(ctx context.Context, commentID, userID uint) (bool, error) {
	return s.repo.HasLiked(ctx, commentID, userID)
}

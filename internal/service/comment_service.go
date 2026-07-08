package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type CommentService interface {
	ListComments(ctx context.Context, chapterID uint, limit, offset int) ([]model.Comment, int64, error)
	ListMyComments(ctx context.Context, userID uint, limit, offset int) ([]model.Comment, int64, error)
	CreateComment(ctx context.Context, chapterID uint, author *model.User, body string) (*model.Comment, error)
	GetComment(ctx context.Context, id uint) (*model.Comment, error)
	DeleteComment(ctx context.Context, id uint) error
}

type commentService struct{ repo repository.CommentRepository }

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (s *commentService) ListComments(ctx context.Context, chapterID uint, limit, offset int) ([]model.Comment, int64, error) {
	list, err := s.repo.ListByChapter(ctx, chapterID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountByChapter(ctx, chapterID)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *commentService) ListMyComments(ctx context.Context, userID uint, limit, offset int) ([]model.Comment, int64, error) {
	list, err := s.repo.ListByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.CountByUser(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// CreateComment takes the authenticated user (already loaded by JWTMiddleware)
// rather than an ID, so the suspension check below reads the state from the
// same request instead of a second DB round-trip.
func (s *commentService) CreateComment(ctx context.Context, chapterID uint, author *model.User, body string) (*model.Comment, error) {
	if !author.HasPermission("comment:create") {
		return nil, errors.New("forbidden")
	}
	if author.CommentSuspendedUntil != nil && author.CommentSuspendedUntil.After(time.Now()) {
		return nil, errors.New("commenting suspended until " + author.CommentSuspendedUntil.Format(time.RFC3339))
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return nil, errors.New("comment body is required")
	}

	c := &model.Comment{ChapterID: chapterID, UserID: author.ID, Body: body}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	c.User = author
	return c, nil
}

func (s *commentService) GetComment(ctx context.Context, id uint) (*model.Comment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *commentService) DeleteComment(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

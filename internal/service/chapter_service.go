package service

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type ChapterService interface {
	GetChapter(ctx context.Context, id uint) (*model.Chapter, error)
	CreateChapter(ctx context.Context, c *model.Chapter) (*model.Chapter, error)
}

type chapterService struct{ repo repository.ChapterRepository }

func NewChapterService(repo repository.ChapterRepository) ChapterService {
	return &chapterService{repo: repo}
}

func (s *chapterService) GetChapter(ctx context.Context, id uint) (*model.Chapter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *chapterService) CreateChapter(ctx context.Context, c *model.Chapter) (*model.Chapter, error) {
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

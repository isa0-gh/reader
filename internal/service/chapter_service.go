package service

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type ChapterService interface {
	GetChapter(ctx context.Context, id uint) (*model.Chapter, error)
	CreateChapter(ctx context.Context, c *model.Chapter) (*model.Chapter, error)
	AddPages(ctx context.Context, chapterID uint, pages []model.S3Object) error
	DeletePage(ctx context.Context, chapterID, pageID uint) error
	DeleteChapter(ctx context.Context, id uint) error
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

func (s *chapterService) AddPages(ctx context.Context, chapterID uint, pages []model.S3Object) error {
	for i := range pages {
		pages[i].ChapterID = &chapterID
	}
	return s.repo.AddPages(ctx, pages)
}

func (s *chapterService) DeletePage(ctx context.Context, chapterID, pageID uint) error {
	return s.repo.DeletePage(ctx, chapterID, pageID)
}

func (s *chapterService) DeleteChapter(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

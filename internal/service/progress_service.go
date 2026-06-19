package service

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type ProgressService interface {
	UpdateProgress(ctx context.Context, userID, chapterID uint, pageIndex int, completed bool) (*model.ReadingProgress, error)
	GetProgress(ctx context.Context, userID, chapterID uint) (*model.ReadingProgress, error)
	GetUserProgress(ctx context.Context, userID uint) ([]model.ReadingProgress, error)
	GetSeriesProgress(ctx context.Context, userID, seriesID uint) ([]model.ReadingProgress, error)
}

type progressService struct{ repo repository.ProgressRepository }

func NewProgressService(repo repository.ProgressRepository) ProgressService {
	return &progressService{repo: repo}
}

func (s *progressService) UpdateProgress(ctx context.Context, userID, chapterID uint, pageIndex int, completed bool) (*model.ReadingProgress, error) {
	p := &model.ReadingProgress{
		UserID:    userID,
		ChapterID: chapterID,
		PageIndex: pageIndex,
		Completed: completed,
	}
	if err := s.repo.Upsert(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *progressService) GetProgress(ctx context.Context, userID, chapterID uint) (*model.ReadingProgress, error) {
	return s.repo.GetByUserAndChapter(ctx, userID, chapterID)
}

func (s *progressService) GetUserProgress(ctx context.Context, userID uint) ([]model.ReadingProgress, error) {
	return s.repo.GetUserProgress(ctx, userID)
}

func (s *progressService) GetSeriesProgress(ctx context.Context, userID, seriesID uint) ([]model.ReadingProgress, error) {
	return s.repo.GetSeriesProgress(ctx, userID, seriesID)
}

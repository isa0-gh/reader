package service

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type SeriesService interface {
	GetSeries(ctx context.Context, id uint) (*model.Series, error)
	CreateSeries(ctx context.Context, s *model.Series) (*model.Series, error)
	UpdateSeries(ctx context.Context, s *model.Series) (*model.Series, error)
	ListSeries(ctx context.Context, limit, offset int, q string) ([]model.Series, int64, error)
	DeleteSeries(ctx context.Context, id uint) error
}

type seriesService struct{ repo repository.SeriesRepository }

func NewSeriesService(repo repository.SeriesRepository) SeriesService {
	return &seriesService{repo: repo}
}

func (s *seriesService) GetSeries(ctx context.Context, id uint) (*model.Series, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *seriesService) CreateSeries(ctx context.Context, series *model.Series) (*model.Series, error) {
	if err := s.repo.Create(ctx, series); err != nil {
		return nil, err
	}
	return series, nil
}

func (s *seriesService) UpdateSeries(ctx context.Context, series *model.Series) (*model.Series, error) {
	if err := s.repo.Update(ctx, series); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, series.ID)
}

func (s *seriesService) ListSeries(ctx context.Context, limit, offset int, q string) ([]model.Series, int64, error) {
	list, err := s.repo.List(ctx, limit, offset, q)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count(ctx, q)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *seriesService) DeleteSeries(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

package service

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type SeriesService interface {
	GetSeries(ctx context.Context, id uint) (*model.Series, error)
	CreateSeries(ctx context.Context, s *model.Series) (*model.Series, error)
	ListSeries(ctx context.Context, query, sort string) ([]model.Series, error)
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

func (s *seriesService) ListSeries(ctx context.Context, query, sort string) ([]model.Series, error) {
	return s.repo.List(ctx, query, sort)
}

func (s *seriesService) DeleteSeries(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

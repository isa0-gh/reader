package service

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type SeriesService interface {
	GetSeries(ctx context.Context, id uint) (*model.Series, error)
}

type seriesService struct{ repo repository.SeriesRepository }

func NewSeriesService(repo repository.SeriesRepository) SeriesService {
	return &seriesService{repo: repo}
}

func (s *seriesService) GetSeries(ctx context.Context, id uint) (*model.Series, error) {
	return s.repo.GetByID(ctx, id)
}

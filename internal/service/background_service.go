package service

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"github.com/isa0-gh/reader/internal/repository"
)

type BackgroundService interface {
	ListBackgrounds(ctx context.Context) ([]model.Background, error)
	AddBackground(ctx context.Context, key, bucket string) (*model.Background, error)
	DeleteBackground(ctx context.Context, id uint) error
}

type backgroundService struct{ repo repository.BackgroundRepository }

func NewBackgroundService(repo repository.BackgroundRepository) BackgroundService {
	return &backgroundService{repo: repo}
}

func (s *backgroundService) ListBackgrounds(ctx context.Context) ([]model.Background, error) {
	return s.repo.List(ctx)
}

func (s *backgroundService) AddBackground(ctx context.Context, key, bucket string) (*model.Background, error) {
	b := &model.Background{Image: &model.S3Object{Key: key, Bucket: bucket}}
	if err := s.repo.Create(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *backgroundService) DeleteBackground(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

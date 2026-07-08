package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type BackgroundRepository interface {
	List(ctx context.Context) ([]model.Background, error)
	Create(ctx context.Context, b *model.Background) error
	Delete(ctx context.Context, id uint) error
}

type backgroundRepository struct{ db *gorm.DB }

func NewBackgroundRepository(db *gorm.DB) BackgroundRepository {
	return &backgroundRepository{db: db}
}

func (r *backgroundRepository) List(ctx context.Context) ([]model.Background, error) {
	// Initialized (not nil) so json.Marshal always encodes "[]", never "null".
	list := []model.Background{}
	if err := r.db.WithContext(ctx).Preload("Image").Order("created_at asc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *backgroundRepository) Create(ctx context.Context, b *model.Background) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(b.Image).Error; err != nil {
			return err
		}
		b.ImageID = b.Image.ID
		return tx.Omit("Image").Create(b).Error
	})
}

func (r *backgroundRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Background{}, id).Error
}

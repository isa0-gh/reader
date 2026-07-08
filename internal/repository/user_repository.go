package repository

import (
	"context"

	"github.com/isa0-gh/reader/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit int, after, before uint) ([]model.User, error)
	Count(ctx context.Context) (int64, error)
	// SetAvatar creates the S3Object row and points the user at it in one
	// transaction, mirroring how seriesRepository.Create/Update handle
	// Series.CoverImage.
	SetAvatar(ctx context.Context, userID uint, avatar *model.S3Object) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Preload("Avatar").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).Preload("Avatar").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update saves the user's own columns only. Omit("Avatar") stops gorm's
// default FullSaveAssociations behavior from re-upserting the preloaded
// Avatar (S3Object) row on every password/role/profile change.
func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Omit("Avatar").Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepository) List(ctx context.Context, limit int, after, before uint) ([]model.User, error) {
	// Initialized (not nil) so json.Marshal always encodes "[]", never
	// "null" — handler.List returns this straight to the frontend, which
	// calls .length on it unconditionally.
	users := []model.User{}
	q := r.db.WithContext(ctx).Order("id asc").Limit(limit)
	if after > 0 {
		q = q.Where("id > ?", after)
	}
	if before > 0 {
		q = q.Where("id < ?", before)
	}
	if err := q.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	return count, r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
}

func (r *userRepository) SetAvatar(ctx context.Context, userID uint, avatar *model.S3Object) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(avatar).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.User{}).Where("id = ?", userID).Update("avatar_id", avatar.ID).Error; err != nil {
			return err
		}
		return tx.Preload("Avatar").First(&user, userID).Error
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

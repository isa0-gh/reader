package model

import "time"

// Background is a site-wide wallpaper image, admin-managed. Leaf table (migrated
// after S3Object), so unlike Series.CoverImage/User.Avatar it can declare a
// normal gorm-managed FK without hitting the migration-order cycle documented
// in database/foreign_keys.go.
type Background struct {
	ID      uint      `gorm:"primaryKey" json:"id"`
	ImageID uint      `gorm:"not null" json:"image_id"`
	Image   *S3Object `gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE" json:"image"`

	CreatedAt time.Time `json:"created_at"`
}

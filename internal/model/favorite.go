package model

import "time"

// Favorite is a leaf table (migrated after User, Series, and Chapter all
// exist), so unlike Series.CoverImage/User.Avatar it can declare normal
// gorm-managed FKs without hitting the migration-order cycle documented in
// database/foreign_keys.go.
type Favorite struct {
	ID     uint  `gorm:"primaryKey" json:"id"`
	UserID uint  `gorm:"not null;uniqueIndex:idx_favorite_user_series" json:"user_id"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`

	SeriesID uint    `gorm:"not null;uniqueIndex:idx_favorite_user_series" json:"series_id"`
	Series   *Series `gorm:"foreignKey:SeriesID;constraint:OnDelete:CASCADE" json:"series,omitempty"`

	LastReadChapterID *uint    `gorm:"index" json:"last_read_chapter_id"`
	LastReadChapter   *Chapter `gorm:"foreignKey:LastReadChapterID;constraint:OnDelete:SET NULL" json:"last_read_chapter,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

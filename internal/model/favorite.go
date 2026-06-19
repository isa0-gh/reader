package model

import "time"

// UserFavorite represents a user's favorite series (many-to-many relationship)
type UserFavorite struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index:idx_user_series,unique"`
	SeriesID  uint      `json:"series_id" gorm:"not null;index:idx_user_series,unique"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"-" gorm:"foreignKey:UserID"`
	Series    *Series   `json:"series,omitempty" gorm:"foreignKey:SeriesID"`
}

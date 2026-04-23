package model

import (
	"time"

	"gorm.io/gorm"
)

type Chapter struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	SeriesID   uint           `gorm:"not null;index" json:"series_id"`
	Series     Series         `gorm:"foreignKey:SeriesID" json:"-"`
	Number     float64        `gorm:"not null" json:"number"`
	Title      string         `json:"title"`
	Pages      []S3Object     `gorm:"foreignKey:ChapterID" json:"pages,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

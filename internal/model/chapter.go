package model

import (
	"time"

	"gorm.io/gorm"
)

type Chapter struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	SeriesID   uint           `gorm:"not null;index;constraint:false" json:"series_id"`
	Series     *Series        `gorm:"-" json:"-"`
	Number     float64        `gorm:"not null" json:"number"`
	Title      string         `json:"title"`
	Pages      []S3Object     `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE" json:"pages,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

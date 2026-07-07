package model

import (
	"time"

	"gorm.io/gorm"
)

type Chapter struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	SeriesID uint `gorm:"not null;index" json:"series_id"`
	// Series is intentionally not a gorm relation (kept out of the schema
	// graph entirely, not just constraint:-) — declaring it would reintroduce
	// the Series->S3Object->Chapter->Series migration cycle CoverImage and
	// Pages create together. The FK itself is added manually in
	// database.EnsureForeignKeys, which runs after AutoMigrate and so isn't
	// exposed to that ordering problem. Repositories load chapters manually.
	Series     *Series        `gorm:"-" json:"-"`
	UploaderID *uint          `gorm:"index" json:"uploader_id,omitempty"`
	Number     float64        `gorm:"not null" json:"number"`
	Title      string         `json:"title"`
	Pages      []S3Object     `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE" json:"pages,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

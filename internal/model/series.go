package model

import (
	"time"

	"gorm.io/gorm"
)

type SeriesStatus string

const (
	StatusOngoing   SeriesStatus = "ongoing"
	StatusCompleted SeriesStatus = "completed"
	StatusHiatus    SeriesStatus = "hiatus"
)

type Series struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Title        string `gorm:"not null;index" json:"title"`
	Slug         string `gorm:"uniqueIndex;not null" json:"slug"`
	Description  string `gorm:"type:text" json:"description"`
	CoverImageID *uint  `gorm:"index" json:"cover_image_id"`
	// constraint:- (not constraint:false, which gorm doesn't recognize and
	// silently ignores) keeps this out of AutoMigrate's own FK management —
	// the actual constraint is created manually in database.EnsureForeignKeys
	// with an explicit ON DELETE SET NULL, since gorm's own auto-created
	// constraint has no ondelete clause and once created can't be altered
	// in place.
	CoverImage *S3Object      `gorm:"foreignKey:CoverImageID;references:ID;constraint:-" json:"cover_image"`
	Author     string         `json:"author"`
	Artist     string         `json:"artist"`
	Status     SeriesStatus   `gorm:"type:varchar(20);default:'ongoing'" json:"status"`
	Chapters   []Chapter      `gorm:"-" json:"chapters,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

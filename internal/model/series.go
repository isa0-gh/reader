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
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null;index" json:"title"`
	Slug        string         `gorm:"uniqueIndex;not null" json:"slug"`
	Description  string         `gorm:"type:text" json:"description"`
	CoverImageID *uint          `gorm:"index" json:"cover_image_id"`
	CoverImage   *S3Object      `gorm:"foreignKey:CoverImageID;references:ID;constraint:false" json:"cover_image"`
	Author       string         `json:"author"`
	Artist       string         `json:"artist"`
	Status       SeriesStatus   `gorm:"type:varchar(20);default:'ongoing'" json:"status"`
	Chapters     []Chapter      `gorm:"-" json:"chapters,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

package model

import (
	"time"
)

type S3Object struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ChapterID  uint      `gorm:"not null;index" json:"chapter_id"`
	Bucket     string    `gorm:"not null" json:"bucket"`
	Key        string    `gorm:"not null;uniqueIndex" json:"key"`
	PageNumber int       `gorm:"not null" json:"page_number"`
	CreatedAt  time.Time `json:"created_at"`
}

package model

import "time"

type ReadingProgress struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index:idx_user_chapter,unique"`
	ChapterID uint      `json:"chapter_id" gorm:"not null;index:idx_user_chapter,unique"`
	PageIndex int       `json:"page_index" gorm:"default:0"`
	Completed bool      `json:"completed" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"-" gorm:"foreignKey:UserID"`
	Chapter   *Chapter  `json:"chapter,omitempty" gorm:"foreignKey:ChapterID"`
}

package model

import "time"

// Comment represents a user comment on a chapter
type Comment struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	ChapterID  uint       `json:"chapter_id" gorm:"not null;index"`
	UserID     uint       `json:"user_id" gorm:"not null;index"`
	ParentID   *uint      `json:"parent_id" gorm:"index"` // For nested replies
	Content    string     `json:"content" gorm:"type:text;not null"`
	Edited     bool       `json:"edited" gorm:"default:false"`
	Deleted    bool       `json:"deleted" gorm:"default:false;index"`
	LikeCount  int        `json:"like_count" gorm:"default:0"`
	ReplyCount int        `json:"reply_count" gorm:"default:0"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	
	Chapter  *Chapter  `json:"chapter,omitempty" gorm:"foreignKey:ChapterID"`
	User     *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Parent   *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Replies  []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
}

// CommentLike represents a user's like on a comment
type CommentLike struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CommentID uint      `json:"comment_id" gorm:"not null;index:idx_comment_user,unique"`
	UserID    uint      `json:"user_id" gorm:"not null;index:idx_comment_user,unique"`
	CreatedAt time.Time `json:"created_at"`
	
	Comment *Comment `json:"-" gorm:"foreignKey:CommentID"`
	User    *User    `json:"-" gorm:"foreignKey:UserID"`
}

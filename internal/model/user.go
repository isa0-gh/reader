package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Name         string `json:"name"`
	Bio          string `gorm:"type:text" json:"bio"`
	AvatarID     *uint  `gorm:"index" json:"avatar_id"`
	// constraint:- for the same reason as Series.CoverImage (see model/series.go):
	// User migrates before S3Object in cmd/server/main.go's AutoMigrate order, so
	// gorm can't create this FK inline. database.EnsureForeignKeys adds it with an
	// explicit ON DELETE SET NULL once every table exists.
	Avatar                *S3Object      `gorm:"foreignKey:AvatarID;references:ID;constraint:-" json:"avatar"`
	Role                  Role           `gorm:"type:varchar(20);default:'reader'" json:"role"`
	JwtID                 string         `gorm:"type:uuid" json:"-"`
	CommentSuspendedUntil *time.Time     `json:"comment_suspended_until,omitempty"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) HasPermission(perm string) bool {
	return u.Role.HasPermission(perm)
}

package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	Email                 string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash          string         `gorm:"not null" json:"-"`
	Name                  string         `json:"name"`
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

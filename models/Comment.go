package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment represents a comment on a post
type Comment struct {
	gorm.Model
	PostID    uint  `json:"post_id" gorm:"foreignKey:PostID"`
	Post      Post  `json:"post"`    // Belongs To relationship with Post
	UserID     uint  `json:"user_id" gorm:"foreignKey:UserID"`
	User       User   `json:"user"`    // Belongs To relationship with User
	Content    string `json:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at"` // Use time.Time for timestamps
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateCommentInput represents data for creating a new comment
type CreateCommentInput struct {
	Content string `json:"content" binding:"required"`
}
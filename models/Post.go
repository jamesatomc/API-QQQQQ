package models

import (
	"time"

	"gorm.io/gorm"
)

// Post represents a post created by a user
type Post struct {
	gorm.Model
	UserID     uint      `json:"user_id" gorm:"foreignKey:UserID"`
	User       User       `json:"user"` // Belongs To relationship with User
	Content    string    `json:"content" validate:"required"`
	CreatedAt time.Time `json:"created_at"` // Use time.Time for timestamps
	UpdatedAt time.Time `json:"updated_at"`
	Comments []Comment `gorm:"foreignKey:PostID"` // Has Many relationship
	Likes    []Like   `gorm:"foreignKey:PostID"` // Has Many relationship
}

// CreatePostInput represents data for creating a new post
type CreatePostInput struct {
	Content string `json:"content" binding:"required"`
}

// UpdatePostInput represents data for updating an existing post
type UpdatePostInput struct {
	Content string `json:"content"`
}
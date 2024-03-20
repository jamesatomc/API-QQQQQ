package models

import (
	"time"

	"gorm.io/gorm"
)


type Like struct {
	gorm.Model
	PostID    uint `json:"post_id" gorm:"foreignKey:PostID"`
	Post      Post  `json:"post"`    // Belongs To relationship with Post
	UserID     uint  `json:"user_id" gorm:"foreignKey:UserID"`
	User       User   `json:"user"`    // Belongs To relationship with User
	CreatedAt time.Time `json:"created_at"` // Use time.Time for timestamps
}

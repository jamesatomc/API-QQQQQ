package models

import (
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	gorm.Model
	Username 	string    	`json:"username" gorm:"unique" validate:"required" `
	Email    	string    	`json:"email" gorm:"unique" validate:"required,email"`
	Password 	string    	`json:"password" validate:"required,min=8"`
	FirstName 	string  	`json:"first_name" validate:"required"`
	LastName  	string 		`json:"last_name" validate:"required"`
	Posts     	[]Post  	`gorm:"foreignKey:UserID"` // Has Many relationship
	Comments 	[]Comment 	`gorm:"foreignKey:UserID"` // Has Many relationship
	Likes    	[]Like   	`gorm:"foreignKey:UserID"` // Has Many relationship
}

// CreateUserInput represents data for creating a new user
type CreateUserInput struct {
	Username 	string 	`json:"username" binding:"required"`
	Email    	string 	`json:"email" binding:"required,email"`
	Password 	string 	`json:"password" binding:"required,min=8"`
	FirstName 	string 	`json:"first_name" binding:"required"`
	LastName  	string 	`json:"last_name" binding:"required"`
}

// UpdateUserInput represents data for updating an existing user
type UpdateUserInput struct {
	Username 	string 	`json:"username"`
	Email    	string 	`json:"email" binding:"email"`
	Password 	string 	`json:"password" binding:"min=8"`
}



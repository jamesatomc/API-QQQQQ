package models
import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username 	string    	  `json:"username" gorm:"unique" validate:"required" `
	Email    	string    	  `json:"email" gorm:"unique" validate:"required,email"`
	Password 	string    	  `json:"password" validate:"required,min=8"`
	KycData 	[]KycData 
	Points      int	  
}

// CreateUserInput represents data for creating a new user
type CreateUserInput struct {
	Username 	string 	`json:"username" binding:"required"`
	Email    	string 	`json:"email" binding:"required,email"`
	Password 	string 	`json:"password" binding:"required,min=8"`
}

// UpdateUserInput represents data for updating an existing user
type UpdateUserInput struct {
	Username 	string 	`json:"username"`
	Email    	string 	`json:"email" binding:"email"`
}

type UpdatePasswordInput struct {
    Username    string `json:"username" binding:"required"`
    OldPassword string `json:"oldpassword" binding:"required"`
    NewPassword string `json:"newpassword" binding:"required,min=8"` // Example: minimum 8 characters 
}


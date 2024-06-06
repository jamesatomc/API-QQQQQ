package models

import (
	"gorm.io/gorm"
)

// In models.go
type KycData struct {
    gorm.Model
    FirstName 	 string 	`json:"firstname" binding:"required"`
    LastName  	 string 	`json:"lastname" binding:"required"`
    IdentityCard string     `json:"identitycard" binding:"required"`
    Country      string    `json:"country" binding:"required"`
    Address      string    `json:"address" binding:"required"`
    IDCardImage  string    `json:"idcardimage" binding:"required"`
    UserID   	 uint
}

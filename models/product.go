package models

import (
	"gorm.io/gorm"
)

type Product struct {
    gorm.Model
    Name        string  `json:"name" binding:"required"`
    Description string  `json:"description" binding:"required"`
    Price       float64 `json:"price" binding:"required"`
}

type UpdateProductInput struct {
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}
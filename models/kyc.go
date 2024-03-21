package models

import "time"

type KYCInput struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UserID             uint      `json:"user_id"` // Link to the associated User
	DocumentType       string    `json:"document_type"`
	DocumentNumber     string    `json:"document_number"`
	DocumentImage      string    `json:"document_image"`                             // Consider storing the path to the image
	VerificationStatus string    `json:"verification_status" gorm:"default:pending"` // pending, approved, rejected
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
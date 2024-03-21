package controllers

import (

	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/jamesatomc/go-api/db"
	"github.com/jamesatomc/go-api/models"

)

func CreateKYC(c *gin.Context) {
    // Get authenticated user ID (e.g., from middleware)
    userID, _ := c.Get("user_id") 

    var input models.KYCInput 
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Handle image upload (if applicable)
    // ...

    kyc := models.KYCInput{
        UserID:         userID.(uint),
        DocumentType:   input.DocumentType,
        DocumentNumber: input.DocumentNumber,
        // DocumentImage:  imagePath, // If you stored an image
    }

    connect.Database.Create(&kyc)
    c.JSON(http.StatusOK, gin.H{"data": kyc})
}
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	connect "github.com/jamesatomc/go-api/db"
	"github.com/jamesatomc/go-api/models"
	"gorm.io/gorm"
)

// In users.go
func AddKycData(c *gin.Context) {
    // Get the username from the request
    username := c.Param("username")

    // Find the user
    var user models.User
    if err := connect.Database.Where("username = ?", username).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
        }
        return
    }

    // Validate input
    var input models.KycData
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if an entry already exists with the same IdentityCard
    var existingEntry models.KycData
    if err := connect.Database.Where("identity_card = ?", input.IdentityCard).First(&existingEntry).Error; err != nil {
        if err != gorm.ErrRecordNotFound {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for existing entry"})
            return
        }
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": "An entry with this identity card already exists"})
        return
    }

    // Add the new entry
    input.UserID = user.ID
    if err := connect.Database.Create(&input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error adding entry"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Entry added successfully", "data": input})
}


func GetKycData(c *gin.Context) {
    // Get the username from the request
    username := c.Param("username")

    // Find the user
    var user models.User
    if err := connect.Database.Where("username = ?", username).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
        }
        return
    }

    // Find the KYC data
    var kycData models.KycData
    if err := connect.Database.Where("user_id = ?", user.ID).First(&kycData).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "KYC data not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding KYC data"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": kycData})
}

func UpdateKycData(c *gin.Context) {
    // Get the username from the request
    username := c.Param("username")

    // Find the user
    var user models.User
    if err := connect.Database.Where("username = ?", username).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
        }
        return
    }

    // Find the KYC data
    var kycData models.KycData
    if err := connect.Database.Where("user_id = ?", user.ID).First(&kycData).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "KYC data not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding KYC data"})
        }
        return
    }

    // Validate input
    var input models.KycData
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Check if an entry already exists with the same IdentityCard
    var existingEntry models.KycData
    if err := connect.Database.Where("identity_card = ?", input.IdentityCard).First(&existingEntry).Error; err != nil {
        if err != gorm.ErrRecordNotFound {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for existing entry"})
            return
        }
    } else {
        if existingEntry.UserID != user.ID {
            c.JSON(http.StatusBadRequest, gin.H{"error": "An entry with this identity card already exists"})
            return
        }
    }

    // Update the KYC data
    if err := connect.Database.Model(&kycData).Updates(input).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating KYC data"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "KYC data updated successfully", "data": kycData})
}
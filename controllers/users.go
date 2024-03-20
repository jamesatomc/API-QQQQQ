package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jamesatomc/go-api/models"
	"gorm.io/gorm"
)


func FindUsers(c *gin.Context) {
	var users []models.User
	models.Database.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func CreateUser(c *gin.Context) {
	// Validate input
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword := sha256.Sum256([]byte(input.Password))
	hashedString := hex.EncodeToString(hashedPassword[:])

	user := models.User{
		Username: input.Username,
		Email:   input.Email,
		Password: hashedString,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	models.Database.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

func FindUser(c *gin.Context) {
	var user models.User

	if err := models.Database.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}



func UpdateUser(c *gin.Context) {
	// Get model if exist
	var user models.User
	if err := models.Database.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
	  return
	}
  
	// Validate input
	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	  return
	}
  
	models.Database.Model(&user).Updates(input)
  
	c.JSON(http.StatusOK, gin.H{"data": user})
}


func DeleteUser(c *gin.Context) {
	// Get model if exist
	var user models.User
	if err := models.Database.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
	  c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
	  return
	}
  
	models.Database.Delete(&user)
  
	c.JSON(http.StatusOK, gin.H{"data": true})
}
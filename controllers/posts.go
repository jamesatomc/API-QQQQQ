package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jamesatomc/go-api/models"
)

func CreatPost(c *gin.Context) {
	// Validate input
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	user := models.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  input.Password,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	models.Database.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}
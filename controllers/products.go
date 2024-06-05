package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	connect "github.com/jamesatomc/go-api/db"
	"github.com/jamesatomc/go-api/models"
)


func CreateProduct(c *gin.Context) {
	var input models.Product

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		Name:		 input.Name,
		Description: input.Description,
		Price:		 input.Price,
	}

	connect.Database.Create(&product)

	c.JSON(http.StatusOK, gin.H{"data": product})
}

func UpdateProduct(c *gin.Context) {
    var product models.Product

    if err := connect.Database.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
        return
    }

    var input models.UpdateProductInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    connect.Database.Model(&product).Updates(models.Product{
        Name:        input.Name,
        Description: input.Description,
        Price:       input.Price,
    })

    c.JSON(http.StatusOK, gin.H{"data": product})
}

func DeleteProduct(c *gin.Context) {
	var product models.Product

	if err := connect.Database.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	connect.Database.Delete(&product)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
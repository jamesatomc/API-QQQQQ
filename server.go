package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jamesatomc/go-api/controllers"
	"github.com/jamesatomc/go-api/models"
	"github.com/joho/godotenv"
)

func main() {

	loadEnv()
    loadDatabase()
    serveApplication()
	
}


func loadEnv() {
    err := godotenv.Load(".env.local")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func loadDatabase() {
    models.ConnectDatabase()
	models.Database.AutoMigrate(&models.User{})
}


func serveApplication() {	

	server := gin.Default()

	models.ConnectDatabase()
	
	// User
	server.GET("/users", controllers.FindUsers)
	server.POST("/register", controllers.CreateUser)
	server.POST("/login", controllers.Login)
	server.GET("/users/:id", controllers.FindUser)
	server.PATCH("/updateusers/:id", controllers.UpdateUser)
	server.PATCH("/change-password/", controllers.UpdatePassword)
	server.DELETE("/users/:id", controllers.DeleteUser)
	

	server.Run()
}


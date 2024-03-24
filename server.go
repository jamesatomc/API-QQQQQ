package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jamesatomc/go-api/controllers"
	"github.com/jamesatomc/go-api/db"
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
    connect.ConnectDatabase()
	connect.Database.AutoMigrate(&models.User{})
}


func serveApplication() {	

	server := gin.Default()

	connect.ConnectDatabase()

	    // User Routes
		userGroup := server.Group("/users")
		{
			// Assuming get all
			userGroup.GET("/", controllers.FindUsers)

			// Assuming get by ID
			userGroup.GET("/:id", controllers.FindUser)

			// Register
			userGroup.POST("/register", controllers.CreateUser)

			// Assuming delete by username
			userGroup.DELETE("/:username", controllers.DeleteUser)

			// Assuming login
			userGroup.POST("/login", controllers.Login)

			// Assuming update by ID
			userGroup.PATCH("/:id", controllers.UpdateUser) 

			// Assuming change password
			userGroup.PATCH("/change-password", controllers.UpdatePassword)
			
			
		}
	

	server.Run()
}


package main

import (
	"log"
	"os"

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

	secretKey := os.Getenv("JWT_SECRET_KEY") // Fetch your secret key

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


			// Routes requiring authentication go within this group
			userGroup.Use(controllers.AuthenticationMiddleware(secretKey)) 
			{
				userGroup.POST("/login", controllers.Login) // Moved inside
				userGroup.PATCH("/:id", controllers.UpdateUser) 
				userGroup.PATCH("/change-password", controllers.UpdatePassword)
			}
			
		}
	

	server.Run()
}


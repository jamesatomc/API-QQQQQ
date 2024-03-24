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
	
	// server.Use(controllers.AuthenticationMiddleware("secret_key"))
	// // User
	// server.GET("/users", controllers.FindUsers)
	// server.GET("/users/:id", controllers.FindUser)
	// server.POST("/register", controllers.CreateUser)
	// server.POST("/login", controllers.Login)
	// server.PATCH("/update-users/:id", controllers.UpdateUser)
	// server.PATCH("/change-password", controllers.UpdatePassword)
    // server.DELETE("/users/:username",  controllers.DeleteUser)

	    // User Routes
		userGroup := server.Group("/users")
		{
			userGroup.GET("/", controllers.FindUsers)
			userGroup.GET("/:id", controllers.FindUser)
			userGroup.POST("/register", controllers.CreateUser)
			userGroup.DELETE("/:username", controllers.DeleteUser)
	
			// Routes requiring authentication go within this group
			secretKey := os.Getenv("SECRET_KEY")
			userGroup.Use(controllers.AuthenticationMiddleware(secretKey))
			{
				userGroup.POST("/login", controllers.Login)
				userGroup.PATCH("/:id", controllers.UpdateUser) // Assuming update by ID
				userGroup.PATCH("/change-password", controllers.UpdatePassword)
			}
		}
	

	server.Run()
}


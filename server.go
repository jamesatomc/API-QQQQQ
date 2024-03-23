package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jamesatomc/go-api/db"
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
    connect.ConnectDatabase()
	connect.Database.AutoMigrate(&models.User{})
}


func serveApplication() {	

	server := gin.Default()

	connect.ConnectDatabase()
	
	// ... Your router setup ...
	authRoutes := server.Group("/")
	authRoutes.Use(controllers.AuthenticationMiddleware()) 
	{
		authRoutes.GET("/users", controllers.FindUsers)
		authRoutes.GET("/users/:id", controllers.FindUser)
		authRoutes.POST("/register", controllers.CreateUser)
		authRoutes.POST("/login", controllers.Login)
		authRoutes.PATCH("/update-users/:id", controllers.UpdateUser)
		authRoutes.PATCH("/change-password", controllers.UpdatePassword)
		authRoutes.DELETE("/users/:username",  controllers.DeleteUser)
	}

	// // User
	// server.GET("/users", controllers.FindUsers)
	// server.GET("/users/:id", controllers.FindUser)
	// server.POST("/register", controllers.CreateUser)
	// server.POST("/login", controllers.Login)
	// server.PATCH("/update-users/:id", controllers.UpdateUser)
	// server.PATCH("/change-password", controllers.UpdatePassword)
    // server.DELETE("/users/:username",  controllers.DeleteUser)
	

	server.Run()
}


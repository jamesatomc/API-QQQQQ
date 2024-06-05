package main

import (
	"log"

	"github.com/gin-gonic/gin"
	middleware "github.com/jamesatomc/go-api/Middleware"
	"github.com/jamesatomc/go-api/controllers"
	"github.com/jamesatomc/go-api/db"
	"github.com/jamesatomc/go-api/models"
	"github.com/joho/godotenv"
)

// main function
func main() {

	loadEnv()
    loadDatabase()
    serveApplication()
	
}

// loadEnv function
func loadEnv() {
    err := godotenv.Load(".env.local")
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

// loadDatabase function
func loadDatabase() {
	// Connect to the database
    connect.ConnectDatabase()
	// Migrate the schema
	connect.Database.AutoMigrate(&models.User{})
	// Migrate the schema
	connect.Database.AutoMigrate(&models.KycData{})
	// Migrate the schema
	connect.Database.AutoMigrate(&models.Product{})
}

// serveApplication function
func serveApplication() {	
	// Create a new server
	server := gin.Default()

	server.Use(middleware.AuthMiddleware())


	// Connect to the database
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

		// KYC Routes
		userKYC := server.Group("/add")
		{	
			// GET KYC 
			userKYC.GET("/kyc/:username", controllers.GetKycData)
			
			// Assuming add KYC data
			userKYC.POST("/kyc/:username", controllers.AddKycData)

			// Assuming update KYC data
			userKYC.PATCH("/kyc/:username", controllers.UpdateKycData)
		}

		// Product Routes
		productGroup := server.Group("/products") 
		{
			// Assuming get all products
			productGroup.POST("/add/products", controllers.CreateProduct)
			
			// Assuming get product by ID
			productGroup.PATCH("/products/:id", controllers.UpdateProduct)
			
			// Assuming delete product by ID
			productGroup.DELETE("/products/:id", controllers.DeleteProduct)
		}

	
		err := server.Run(":8080")
		if err != nil {
		   panic(err)
		}

	
	// server.Run()
}
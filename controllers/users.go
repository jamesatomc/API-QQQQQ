package controllers

// Import necessary packages
import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jamesatomc/go-api/db"
	"github.com/jamesatomc/go-api/models"
	"gorm.io/gorm"
)

// Argon2 Hashing Function
func hashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// Argon2 Password Comparison Function
func comparePassword(hashedPassword, password string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hashedPassword)
	return err == nil && match
}

// FindUsers function
func FindUsers(c *gin.Context) {
	var users []models.User
	connect.Database.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// CreateUser function
func CreateUser(c *gin.Context) {
	// Validate input
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	user := models.User{
		Username:  input.Username,
		Email:     input.Email,
		Password:  hashedPassword,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	// Check for duplicate username
	var existingUser models.User
	if err := connect.Database.Where("username = ?", input.Username).First(&existingUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			// Handle other database errors
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for username"})
			return
		}
		// Record not found - username is available
	} else {
		// Username already exists
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}
	// Use result instead of directly saving
	result := connect.Database.Create(&user)
	// Error handling:
	if result.Error != nil {
		// Check if the error is due to a duplicate email
		if strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		}
		return
	}

	connect.Database.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// FindUser function
func FindUser(c *gin.Context) {
	var user models.User

	if err := connect.Database.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser function
func UpdateUser(c *gin.Context) {
	// Get model if exist
	var user models.User
	if err := connect.Database.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
		}
		return
	}

	// Validate input
	var input models.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate email (if the email is being changed)
	if input.Email != user.Email { // Only check if there's a change
		var existingUser models.User
		if err := connect.Database.Where("email = ?", input.Email).First(&existingUser).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				// Handle other database errors
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for email"})
				return
			}
			// else -> Record not found, so email is available
		} else {
			// Email already exists
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
	}

	// Check for duplicate username (if the username is being changed)
	if input.Username != user.Username {
		var existingUser models.User
		if err := connect.Database.Where("username = ?", input.Username).First(&existingUser).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking for username"})
				return
			}
			// else -> Record not found, so username is available
		} else {
			// Username already exists
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
	}

	// Update the user record
	connect.Database.Model(&user).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// Login function
func Login(c *gin.Context) {
	var input models.User // Use a specific struct for login credentials

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := connect.Database.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Avoid revealing if it's username or password issue
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		// Handle other potential database errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging in"})
		return
	}

	// Compare hashed password
	if !comparePassword(user.Password, input.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate authentication token (consider using JWT)
	token, err := GenerateToken(user.ID, time.Hour*24*7) // Token valid for 1 week
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	// Set the token in a secure, HttpOnly cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 7), // 1 week from now
		HttpOnly: true,                               // Prevent JavaScript access
		Secure:   true,                               // Enforce HTTPS transmission
		Path:     "/",                                // Apply cookie to all paths within the domain
		SameSite: http.SameSiteStrictMode,            // Prevent CSRF
	})

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"}) // Don't send full token back
}

// GenerateToken function
func GenerateToken(userID uint, expiration time.Duration) (string, error) {
	// Move secret key retrieval and storage outside the function (refer to previous improvements)
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY environment variable not set")
	}

	token := jwt.New(jwt.SigningMethodRS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(expiration).Unix() // Set expiration

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// DeleteUser function
func DeleteUser(c *gin.Context) {
	// Get the username from the request (route parameters, query, etc.)
	username := c.Param("username") // Example: Assuming username in route parameter

	// Find the user to delete
	var user models.User
	if err := connect.Database.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		}
		return
	}

	// Delete the user
	if err := connect.Database.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// UpdatePassword function
func UpdatePassword(c *gin.Context) {
	var input models.UpdatePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Find the user by username
	var user models.User
	if err := connect.Database.Where("username = ?", input.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding user"})
		}
		return
	}

	// 2. Verify old password
	if !comparePassword(user.Password, input.OldPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect old password"})
		return
	}

	// 3. Hash the new password
	newPasswordHashed, err := hashPassword(input.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing new password"})
		return
	}

	// 4. Update the user's password
	connect.Database.Model(&user).Update("password", newPasswordHashed)

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

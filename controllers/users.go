package controllers

import (
	// "crypto/sha256"
	// "encoding/hex"
    
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	// "golang.org/x/crypto/bcrypt"

    "github.com/alexedwards/argon2id"

	"github.com/jamesatomc/go-api/db"
	"github.com/jamesatomc/go-api/models"
	"gorm.io/gorm"
)


// func hashPassword(password string) (string, error) {
//     hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Adjust cost as needed
//     if err != nil {
//       return "", err
//     }
//     return string(hash), nil
// }
  
// func comparePassword(hashedPassword, password string) bool {
//     err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
//     return err == nil
// }

func FindUsers(c *gin.Context) {
	var users []models.User
	connect.Database.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func CreateUser(c *gin.Context) {
	// Validate input
	var input models.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hashedPassword := sha256.Sum256([]byte(input.Password))
	// hashedString := hex.EncodeToString(hashedPassword[:])

    // hashedPassword, err := hashPassword(input.Password)
    // if err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
    //     return
    // }

    hashedPassword, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }
    
	user := models.User{
		Username: input.Username,
		Email:   input.Email,
		Password: hashedPassword,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	connect.Database.Create(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

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



func UpdateUser(c *gin.Context) {
    // Get model if exist
    var user models.User
    if err := connect.Database.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
      return
    }
  
    // Validate input
    var input models.UpdateUserInput
    if err := c.ShouldBindJSON(&input); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }
  
    connect.Database.Model(&user).Updates(input)
  
    c.JSON(http.StatusOK, gin.H{"data": user})
}


func Login(c *gin.Context) {
    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := connect.Database.Where("username = ?", input.Username).First(&user).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
            return
        }
        // Handle other potential database errors
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error logging in"})
        return
    }

    // Compare hashed password
    // inputHash := sha256.Sum256([]byte(input.Password))
    // if hex.EncodeToString(inputHash[:]) != user.Password {
    //     c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
    //     return
    // }

    // if !comparePassword(user.Password, input.Password) {
    //     c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
    //     return
    // }

    match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying password"})
        return 
    }
    if !match {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
        return
    } 

    // Generate authentication token (consider using JWT)
    token, err := GenerateToken(user.ID) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
        return 
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

func DeleteUser(c *gin.Context) {
    // Get the username from the request (route parameters, query, etc.)
    username := c.Param("username")  // Example: Assuming username in route parameter

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

func GenerateToken(userID uint) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = userID
    claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Example expiration

   tokenString, err := token.SignedString([]byte("your_strong_secret_key"))
   if err != nil {
       return "", err
   }

   return tokenString, nil
}

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
    // oldPasswordHash := sha256.Sum256([]byte(input.OldPassword))
    // if hex.EncodeToString(oldPasswordHash[:]) != user.Password {
    //     c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect old password"})
    //     return
    // }

    // if !comparePassword(user.Password, input.OldPassword) {
    //     c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect old password"})
    //     return
    // }


    // 2. Verify old password using Argon2
    match, err := argon2id.ComparePasswordAndHash(input.OldPassword, user.Password) 
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verifying old password"})
        return
    }
    if !match {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect old password"})
        return
    }

    // 3. Hash the new password
    // newPasswordHash := sha256.Sum256([]byte(input.NewPassword))
    // newPasswordHashedString := hex.EncodeToString(newPasswordHash[:])

    // newPasswordHashed, err := hashPassword(input.NewPassword) 
    // if err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing new password"})
    //     return
    // }


  // 3. Hash the new password using Argon2
    newPasswordHashed, err := argon2id.CreateHash(input.NewPassword, argon2id.DefaultParams)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing new password"})
        return
    }

    // 4. Update the user's password
    // connect.Database.Model(&user).Update("password", newPasswordHashedString)
    connect.Database.Model(&user).Update("password", newPasswordHashed) 

    c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
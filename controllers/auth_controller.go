package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
)

// AuthenticationMiddleware function
func AuthenticationMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Retrieve token from cookie
        tokenString, err := c.Cookie("auth_token")
        if err != nil || tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        // Validate token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(os.Getenv("JWT_SECRET_KEY")), nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        userID := uint(claims["user_id"].(float64)) 
        c.Set("user_id", userID) // Set the user ID in the Gin context

        c.Next() 
    }
}

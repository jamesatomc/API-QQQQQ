package controllers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
    jwt.StandardClaims
    UserID uint `json:"user_id"`
}

func AuthenticationMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Request.Cookie("auth_token")
		if err != nil || tokenString.Value == "" { // Check for token presence
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing token"})
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
				return
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error parsing token"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			return
		}

		// Add user_id to context
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

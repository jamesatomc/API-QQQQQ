package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)


func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
		secretKey := os.Getenv("JWT_SECRET_KEY")
        authHeader := c.GetHeader("Authorization")
        bearerToken := strings.Split(authHeader, " ")

        if len(bearerToken) != 2 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        token := bearerToken[1]
        // Validate the token here. This is just a placeholder.
        if token != secretKey {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        c.Next()
    }
}




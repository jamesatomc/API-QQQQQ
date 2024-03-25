package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)
const secretKey = "fuiwufiquth3uityh3tnc3n23un32ut3667#ys?kxs"

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }

        // Split into "Bearer <token>"
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization format"})
            c.Abort()
            return
        }

        tokenString := parts[1]

        // Validate JWT token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // Verify the signing algorithm and secret/key
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(secretKey), nil 
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            userID := claims["user_id"] // Example of extracting data
            c.Set("user_id", userID)    // Store in Gin context        
        }
        // ... (rest of your middleware code)
        c.Next()
        // ...
    }
}

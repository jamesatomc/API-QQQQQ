package middleware

import (
    "errors"
    "net/http"
    "os"
    "strings"
  

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt"
)

// CustomClaims to add additional data to the token
type CustomClaims struct {
    UserID uint `json:"user_id"`
    jwt.StandardClaims
}

// Authorization Middleware
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString, err := extractToken(c.Request)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        token, err := validateToken(tokenString)
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        // Pass user information down the chain
        claims := token.Claims.(CustomClaims)
        c.Set("user_id", claims.UserID)

        c.Next() // Continue processing the request
    }
}

// Helper functions to extract and validate the token
func extractToken(r *http.Request) (string, error) {
    header := r.Header.Get("Authorization")
    if header == "" {
        return "", errors.New("Missing Authorization header")
    }
    parts := strings.Split(header, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return "", errors.New("Invalid Authorization header format")
    }
    return parts[1], nil
}

func validateToken(tokenString string) (*jwt.Token, error) {
    secretKey := os.Getenv("JWT_SECRET_KEY")
    token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })
    return token, err
}

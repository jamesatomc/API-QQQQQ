package controllers

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)


type Claims struct {
    jwt.StandardClaims
    UserID uint
  }
  

func AuthenticationMiddleware(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
      token, err := c.Request.Cookie("auth_token")
      if err != nil {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
      }
  
      claims := &Claims{}
      parsedToken, err := jwt.ParseWithClaims(token.Value, claims, func(token *jwt.Token) (interface{}, error) {
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
  
      if !parsedToken.Valid {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
        return
      }
  
      // เพิ่ม user_id ลงใน context
      c.Set("user_id", claims.UserID)
  
      c.Next()
    }
  }
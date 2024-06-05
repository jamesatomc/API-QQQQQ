package helper

import (
	"os"

	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken function
func GenerateToken(userID uint, expiration time.Duration) (string, error) {
    // Move secret key retrieval and storage outside the function (refer to previous improvements)
    secretKey := os.Getenv("JWT_SECRET_KEY")

    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = userID
    claims["exp"] = time.Now().Add(expiration).Unix() // Set expiration

    tokenString, err := token.SignedString([]byte(secretKey))
    if err != nil {
      return "", err
    }

    return tokenString, nil
}



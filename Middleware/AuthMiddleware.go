package middleware

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func AuthMiddleware() gin.HandlerFunc {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Erro ao ler variaveis de ambiente")
    }
    requiredToken := os.Getenv("API_TOKEN")

    if requiredToken == "" {
        log.Fatal("Por favor, defina a variavel API_TOKEN")
    }

    return func(c *gin.Context) {
        token := c.Request.FormValue("api_token")

        if token == "" {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Token deve ser preenchido"})
            return
        }
        if token != requiredToken {
            c.JSON(http.StatusBadRequest, gin.H{"message": "Token invalido"})
            return
        }

        c.Next()
    }
}
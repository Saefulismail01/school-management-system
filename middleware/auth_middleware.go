package middleware

import (
	// "fmt"
	"net/http"
	"os"
	"strings"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	// "bimbel-absensi/config"
	// "bimbel-absensi/models"
)

// Auth middleware untuk memeriksa token JWT
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := jwt.MapClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user info in context
		userID := uint(claims["user_id"].(float64))
		role := claims["role"].(string)
		
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Next()
	}
}
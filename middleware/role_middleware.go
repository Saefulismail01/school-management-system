package middleware

import (
	// "fmt"
	"net/http"
	"os"
	// "strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	// "bimbel-absensi/config"
	"bimbel-absensi/models"
)

// RoleAdmin middleware untuk memeriksa apakah user adalah admin
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RolePengajar middleware untuk memeriksa apakah user adalah pengajar
func PengajarOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "pengajar" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Pengajar access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GenerateToken untuk membuat JWT token
func GenerateToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}
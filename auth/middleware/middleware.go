package middleware

import (
	"auth-service/auth/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		claims, err := token.ExtractClaim(authHeader)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		if claims.ExpiresAt < time.Now().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			c.Abort()
			return
		}

		// Claimsdan ma'lumotlarni kontekstga qo'shish
		c.Set("user_id", claims.UserId)
		c.Set("username", claims.Username)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

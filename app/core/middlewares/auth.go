package middlewares

import (
	"ecommerce/app/models"
	"net/http"
	"strings"

	"ecommerce/app/core"
	"ecommerce/app/core/security"
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

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		accessToken := bearerToken[1]
		db := core.GetDB()

		claims, err := security.ValidateToken(db, accessToken, "access")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}
		var user models.User

		result := db.Where("email = ?", claims.Email).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		// Set the authenticated user's email in the context
		c.Set("user", user)

		c.Next()
	}
}

package middlewares

import (
	"net/http"

	"golang-phone-review-api/utils/token"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
				c.Abort()
				return
		}

		err := token.TokenValid(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token."})
			c.Abort()
			return
		}
		c.Next()
	}
}
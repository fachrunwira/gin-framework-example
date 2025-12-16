package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": false,
				"errors": "authorization_expected",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

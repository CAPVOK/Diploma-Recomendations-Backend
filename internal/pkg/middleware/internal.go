package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Internal(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := c.GetHeader("X-Internal-Token"); token != expectedToken {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}

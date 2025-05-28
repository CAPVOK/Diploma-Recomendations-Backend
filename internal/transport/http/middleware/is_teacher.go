package middleware

import (
	"diprec_api/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

func OnlyTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "" || role == domain.RoleStudent.String() {
			c.AbortWithStatusJSON(http.StatusForbidden, domain.Error{Message: domain.ErrInvalidRole.Error()})
			return
		}

		c.Next()
	}
}

package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func (m *middleware) Origin(c *gin.Context) {
	origin := c.GetHeader("Origin")
	if origin == "" {
		c.Next()
		return
	}

	if !slices.Contains(m.allowedOrigins, origin) {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": fmt.Sprintf("request %s not allowed", origin),
		})
		c.Abort()
		return
	}
	c.Next()
}

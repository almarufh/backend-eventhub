package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) Origin(c *gin.Context) {
	origin := c.GetHeader("Origin")
	if origin == "" {
		c.Next()
		// c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		// 	"status":  false,
		// 	"message": "request must have origin allowed",
		// })
		return
	}

	if !slices.Contains(m.allowedOrigins, origin) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": fmt.Sprintf("request %s not allowed", origin),
		})
		return
	}
	c.Next()
}

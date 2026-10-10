package middleware

import (
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) Origin(c *gin.Context) {
	origin := c.GetHeader("Origin")
	log.Printf("Origin Request : \n\n%s\n\n", origin)

	if origin == "" {
		// c.Next()
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": "request must have origin allowed",
		})
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

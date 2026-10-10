package middleware

import (
	"fmt"
	"log"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) Origin(c *gin.Context) {
	client, _ := m.GetClientInfo(c)

	log.Println("Request Accepted from : ")
	fmt.Println("    IP         : ", client.IP)
	fmt.Println("    Origin     : ", client.Origin)
	fmt.Println("    User Agent : ", client.UserAgent)

	if client.Origin == "" {
		// c.Next()
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": "request must have origin allowed",
		})
		return
	}

	if !slices.Contains(m.allowedOrigins, client.Origin) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"status":  false,
			"message": fmt.Sprintf("request %s not allowed", client.Origin),
		})
		return
	}
	c.Next()
}

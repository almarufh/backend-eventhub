package config

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IPWhitelistMiddleware(allowedIPs []string) gin.HandlerFunc {
	ipMap := make(map[string]bool)
	for _, ip := range allowedIPs {
		ipMap[ip] = true
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		if !ipMap[clientIP] {
			log.Printf("[IP Whitelist] Blocked unauthorized IP: %s\n", clientIP)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access forbidden: IP not allowed",
			})
			return
		}

		c.Next()
	}
}

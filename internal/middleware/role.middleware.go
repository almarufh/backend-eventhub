package middleware

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/pkg"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) RequireAdmin(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		log.Println("[RequireAdmin] Error: token context not found in request")
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Authentication required",
		})
		return
	}

	payload, ok := token.(pkg.JWTClaims)
	if !ok {
		log.Printf("[RequireAdmin] Error: failed to cast token to pkg.JWTClaims (type: %T)\n", token)
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Invalid token claims",
		})
		return
	}

	if payload.Role != "admin" {
		log.Printf("[RequireAdmin] Forbidden: user ID %d has role '%s', expected 'admin'\n", payload.ID, payload.Role)
		c.AbortWithStatusJSON(http.StatusForbidden, dto.Response{
			Success: false,
			Message: "Admin access required",
		})
		return
	}

	c.Next()
}

func (m *Middleware) RequireOrganizer(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		log.Println("[RequireOrganizer] Error: token context not found in request")
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Authentication required",
		})
		return
	}

	payload, ok := token.(pkg.JWTClaims)
	if !ok {
		log.Printf("[RequireOrganizer] Error: failed to cast token to pkg.JWTClaims (type: %T)\n", token)
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Invalid token claims",
		})
		return
	}

	if payload.Role != "organizer" {
		log.Printf("[RequireOrganizer] Forbidden: user ID %d has role '%s', expected 'organizer'\n", payload.ID, payload.Role)
		c.AbortWithStatusJSON(http.StatusForbidden, dto.Response{
			Success: false,
			Message: "Organizer access required",
		})
		return
	}

	c.Next()
}

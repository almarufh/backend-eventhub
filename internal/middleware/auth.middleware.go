package middleware

import (
	"backend/EventHub/internal/dto"
	"backend/EventHub/pkg"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) AuthMiddleware(c *gin.Context) {
	log.Println("Request from : ", c.ClientIP())
	payload, err := m.GetPayload(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	authSession, err := m.authRepo.GetAuthUser(c.Request.Context(), int32(payload.ID))
	if err != nil {
		log.Printf("[AuthMiddleware] Session not found for id %d: %v\n", payload.ID, err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Session invalid or expired",
		})
		return
	}

	if !authSession.IsActive {
		log.Printf("[AuthMiddleware] Session %d is revoked/logged out\n", payload.ID)
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Session has ended. Please log in again.",
		})
		return
	}

	c.Set("payload", payload)
	c.Next()
}

func (m *Middleware) CekDong(c *gin.Context) {
	var token pkg.JWTClaims
	err := token.DecodeToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo0LCJyb2xlIjoiYWRtaW4iLCJpc3MiOiJhbG1hcnVmaGlkYXlhdCIsImV4cCI6MTc5MDQxODU1NX0.eQloYcRDFFFecI4Py_VBz3cPuE__9DyErf0R3G9pFGI")
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			log.Printf("[GET PAYLOAD] Warning: token validation failed: %v\n", err)
		}

		log.Printf("[GET PAYLOAD] Critical Error: unexpected error while decoding token: %v\n", err)
	}

	var expiresInSec int64 = 0
	if token.ExpiresAt != nil {
		diff := time.Until(token.ExpiresAt.Time)
		if diff > 0 {
			expiresInSec = int64(diff.Seconds())
		}
	}

	log.Printf("ID      : %d\n", token.ID)
	log.Printf("Role    : %s\n", token.Role)
	log.Printf("Expired : %d\n", expiresInSec)
}

package middleware

import (
	"backend/EventHub/internal/dto"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

	key := fmt.Sprintf("Token:%d", payload.ID)
	res, err := m.redis.Get(c.Request.Context(), key).Result()
	if err == nil {
		log.Printf("[Redis.Logout] Session for user %d is revoked/logged out (Token: %s)\n", payload.ID, res)
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Session has ended. Please log in again.",
		})
		return
	}

	authSession, err := m.authRepo.GetAuthUser(c.Request.Context(), m.pool, int32(payload.ID))
	if err != nil {
		log.Printf("[AuthMiddleware] Session not found for id %d: %v\n", payload.ID, err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Session invalid or expired",
		})
		return
	}

	if !authSession.IsActive {
		log.Printf("[AuthMiddleware] Session %d is inactive/revoked in DB\n", payload.ID)
		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Session has ended. Please log in again.",
		})
		return
	}

	c.Set("payload", payload)
	c.Next()
}

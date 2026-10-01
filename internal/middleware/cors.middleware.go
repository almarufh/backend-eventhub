package middleware

import (
	"backend/EventHub/internal/repo"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Payload struct {
	ID        int32
	Role      string
	ExpiresIn int64
}

type Middleware struct {
	allowedOrigins []string
	authRepo       *repo.AuthRepo
}

func InitMiddleWare(origins []string) *Middleware {
	return &Middleware{
		allowedOrigins: origins,
	}
}

func (m *Middleware) RepoAuthMiddleWare(authRepo *repo.AuthRepo) {
	m.authRepo = authRepo
}

func (m *Middleware) Cors(c *gin.Context) {
	origin := c.GetHeader("Origin")
	c.Header("Access-Control-Allow-Origin", origin)

	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, PATCH")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

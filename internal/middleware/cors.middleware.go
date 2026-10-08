package middleware

import (
	"backend/EventHub/internal/repo"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Payload struct {
	ID        int32
	Role      string
	ExpiresIn int64
}

type Middleware struct {
	allowedOrigins []string
	authRepo       *repo.AuthRepo
	pool           *pgxpool.Pool
	redis          *redis.Client
}

func InitMiddleWare(origins []string, pool *pgxpool.Pool, redis *redis.Client) *Middleware {
	return &Middleware{
		allowedOrigins: origins,
		pool:           pool,
		redis:          redis,
	}
}

func (m *Middleware) RepoAuthMiddleWare(authRepo *repo.AuthRepo) {
	m.authRepo = authRepo
}

func (m *Middleware) Cors(c *gin.Context) {
	origin := c.GetHeader("Origin")
	log.Printf("Origin : %s\n\n", origin)
	c.Header("Access-Control-Allow-Origin", origin)

	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Methods", "DELETE, OPTIONS, PATCH")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

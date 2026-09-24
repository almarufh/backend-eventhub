package router

import (
	"backend/EventHub/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	router     *gin.Engine
	pool       *pgxpool.Pool
	middleware *middleware.Middleware
}

func NewRouter(router *gin.Engine, pool *pgxpool.Pool, middleware *middleware.Middleware) *Router {
	return &Router{
		router:     router,
		pool:       pool,
		middleware: middleware,
	}
}

func (r *Router) Connect() {
	r.statusServer()
	r.router.Use(r.middleware.Origin, r.middleware.Cors)
	r.usersRouter()
}

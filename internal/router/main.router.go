package router

import (
	"backend/EventHub/internal/middleware"

	_ "backend/EventHub/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	router     *gin.Engine
	pool       *pgxpool.Pool
	middleware *middleware.Middleware
	redis      *redis.Client
}

func NewRouter(router *gin.Engine, pool *pgxpool.Pool, middleware *middleware.Middleware, redis *redis.Client) *Router {
	return &Router{
		router:     router,
		pool:       pool,
		middleware: middleware,
		redis:      redis,
	}
}

func (r *Router) Connect() {

	r.router.GET("/docs/v1/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.statusServer()
	r.router.Use(r.middleware.Origin, r.middleware.Cors)
	r.authRouter()
	r.userRouter()
	r.communitieRouter()
}

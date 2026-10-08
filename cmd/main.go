package main

import (
	"backend/EventHub/internal/config"
	"backend/EventHub/internal/middleware"
	"backend/EventHub/internal/router"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

// @title           			Application EventHub
// @version         			1.0
// @description     			This is a sample send request use API.

// @host      					192.168.0.186:7202
// @BasePath  					/

// @securityDefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @description					Bearer Token used as identity for accessing backend
func main() {
	// ENV
	conf := config.LoadEnv()

	// Redis
	redis := conf.REDIS.Connect()
	defer func() {
		if err := redis.Close(); err != nil {
			log.Println("Close Redis : ", err.Error())
		}
	}()

	// PostgreSQL
	pool := conf.POSTGRES.Connect()
	defer pool.Close()

	// Gin Gonic
	server := gin.Default()

	log.Printf("Allow Origins : %s\n\n", conf.ALLOWED_ORIGINS)
	middleware := middleware.InitMiddleWare(conf.ALLOWED_ORIGINS, pool, redis)
	router.NewRouter(server, pool, middleware, redis).Connect()
	server.Run(fmt.Sprintf("%s:%s", conf.HOST, conf.PORT))
}

package main

import (
	"backend/EventHub/internal/router"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	server := gin.Default()
	router.InitMainRouter(server)

	server.Run(fmt.Sprintf("%s:%s", host, port))
}

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	r := gin.Default()
	r.GET("/status", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.String(http.StatusOK, "<h1>Server is running!</h1>")
	})
	r.Run(fmt.Sprintf("%s:%s", host, port))
}

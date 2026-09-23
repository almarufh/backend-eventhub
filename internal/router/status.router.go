package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func statusServer(r *gin.Engine) {
	r.GET("/status", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.String(http.StatusOK, "<h1>Server is running!</h1>")
	})
}

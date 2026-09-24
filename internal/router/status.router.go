package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <style>
        body {
            margin: 0;
            height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            background-color: #000000;
        }
        h1 {
            color: #00ff00;
            font-family: sans-serif;
        }
    </style>
</head>
<body>
    <h1>Server is running!</h1>
</body>
</html>
`

func (r *Router) statusServer() {
	r.router.GET("", func(ctx *gin.Context) {
		ctx.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	})
}

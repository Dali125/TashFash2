package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
)

func Logging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println(ctx.Request.Method, ctx.Request.URL.Path)
		ctx.Next()
	}
}

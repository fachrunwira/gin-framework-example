package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetLog(filename string) gin.HandlerFunc {
	lomberjackLog := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	logger := log.New(lomberjackLog, "", log.LstdFlags)

	return func(ctx *gin.Context) {
		start := time.Now()

		ctx.Next()

		// Log details
		latency := time.Since(start)
		clientIP := ctx.ClientIP()
		method := ctx.Request.Method
		statusCode := ctx.Writer.Status()
		path := ctx.Request.URL.Path
		userAgent := ctx.Request.UserAgent()

		logger.Printf("[GIN] %v | %3d | %13v | %15s | %-7s %s | %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			userAgent,
		)
	}
}

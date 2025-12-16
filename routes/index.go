package routes

import (
	"log/slog"

	"github.com/fachrunwira/gin-example/lib/logger"
	"github.com/gin-gonic/gin"
)

var appLogger *slog.Logger = logger.New("./storage/logs/app.log", slog.LevelInfo)

func RegisterRoutes(g *gin.Engine) {
	userRoutes(g)
}

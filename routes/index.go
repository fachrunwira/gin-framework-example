package routes

import (
	"log/slog"

	"github.com/fachrunwira/gin-example/app/response"
	"github.com/fachrunwira/gin-example/lib"
	"github.com/gin-gonic/gin"
)

var appLogger *slog.Logger = lib.NewLogger("./storage/logs/app.log", slog.LevelInfo)

func RegisterRoutes(g *gin.Engine) {
	userRoutes(g)

	g.NoRoute(func(ctx *gin.Context) {
		response.RouteNotFound(ctx)
	})
}

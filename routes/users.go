package routes

import (
	"github.com/fachrunwira/gin-example/controllers/users"
	"github.com/gin-gonic/gin"
)

func userRoutes(g *gin.Engine) {
	userControllers := users.UserControllers(appLogger)
	userGroup := g.Group("/user")
	userGroup.GET("/list", userControllers.List)
}

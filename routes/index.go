package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.Engine) {
	userRoutes(g)
}

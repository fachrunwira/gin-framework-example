package users

import (
	"log/slog"

	"github.com/fachrunwira/gin-example/app/response"
	"github.com/gin-gonic/gin"
)

type structController struct {
	logger *slog.Logger
}

func UserControllers(logger *slog.Logger) *structController {
	return &structController{
		logger: logger,
	}
}

func (sc *structController) List(c *gin.Context) {
	var singleUser fetchSingleUserDTO

	if err := c.ShouldBindQuery(&singleUser); err != nil {
		sc.logger.Error("failed to bind query params", "err", err)
		c.JSON(500, gin.H{
			"status":   false,
			"message":  "internal server error",
			"contents": nil,
		})
		return
	}

	response.Ok(c, "Ok", true)
}

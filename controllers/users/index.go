package users

import (
	"database/sql"
	"errors"

	"github.com/fachrunwira/go-query-builder/builder"
	"github.com/gin-gonic/gin"
)

type fetchSingleUserDTO struct {
	ID int `form:"id"`
}

func (uc *userControllers) List(c *gin.Context) {
	ctx := c.Request.Context()
	var singleUser fetchSingleUserDTO

	if err := c.ShouldBindQuery(&singleUser); err != nil {
		uc.appLog.Error("failed to bind query params", "err", err)
		c.JSON(500, gin.H{
			"status":   false,
			"message":  "internal server error",
			"contents": nil,
		})
		return
	}

	uc.appLog.Info("user", "id", singleUser.ID)

	result, err := builder.
		MakeWithContext(ctx).
		Table("users").
		Where("id", singleUser.ID).
		First()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(200, gin.H{
				"status":   true,
				"message":  "user not found",
				"contents": result,
			})
			return
		}

		c.JSON(500, gin.H{
			"status":   false,
			"message":  "internal server error",
			"contents": nil,
		})
		uc.appLog.Error("Terjadi kesalahan", "error", err)
		return
	}

	c.JSON(200, gin.H{
		"status":   true,
		"message":  "user found",
		"contents": result,
	})
}

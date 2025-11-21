package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/fachrunwira/ebookamd-api/database"
	"github.com/fachrunwira/go-query-builder/builder"
	"github.com/gin-gonic/gin"
)

func userRoutes(g *gin.Engine) {
	user := g.Group("/user")
	user.GET("/list", func(c *gin.Context) {
		db := database.GetDB()

		result, err := builder.Make(db).Table("author").Select("id", "name").First()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusOK, gin.H{
					"status":  true,
					"content": result,
					"message": "user not found",
				})
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"content": nil,
				"message": "internal server error",
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"content": result,
			"message": "user found",
		})
	})
}

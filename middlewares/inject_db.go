package middlewares

import (
	"github.com/fachrunwira/ebookamd-api/database"
	"github.com/gin-gonic/gin"
)

func InjectDB() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbCtx := database.Inject(ctx.Request.Context())
		ctx.Request = ctx.Request.WithContext(dbCtx)

		ctx.Next()
	}
}

package middlewares

import (
	"errors"
	"strings"

	"github.com/fachrunwira/gin-example/app/response"
	"github.com/fachrunwira/gin-example/lib"
	"github.com/fachrunwira/gin-example/lib/env"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthAccess(includeExpired bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			response.Unauthorized(c, "Authorization header expected")
			c.Abort()
			return
		}

		headers := strings.Split(header, " ")
		if len(headers) != 2 {
			response.Unauthorized(c, "Incorrect authorization format")
			c.Abort()
			return
		}

		if strings.ToLower(headers[0]) != "bearer" {
			response.Unauthorized(c, "Incorrect authorization type")
			c.Abort()
			return
		}

		claims, err := lib.ValidateToken(headers[1])
		if err != nil {
			tokenErrors(c, includeExpired, err)
			c.Abort()
			return
		}

		iss, ok := claims["iss"].(string)
		if !ok {
			response.Unauthorized(c, "Malformed issuer")
			c.Abort()
			return
		}

		if iss != env.Get("APP_URL", "http://localhost") {
			response.Unauthorized(c, "Invalid Issuer")
			c.Abort()
			return
		}

		c.Next()
	}
}

func tokenErrors(c *gin.Context, flag bool, err error) {
	if flag {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			response.Unauthorized(c, "Token Expired")
		case errors.Is(err, jwt.ErrInvalidKey):
			response.Unauthorized(c, "Invalid authorization")
		case errors.Is(err, jwt.ErrSignatureInvalid):
			response.Unauthorized(c, "Invalid signature")
		default:
			response.Unauthorized(c, "Unknown token")
		}
		return
	}

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
	case errors.Is(err, jwt.ErrInvalidKey):
		response.Unauthorized(c, "Invalid authorization")
	case errors.Is(err, jwt.ErrSignatureInvalid):
		response.Unauthorized(c, "Invalid signature")
	default:
		response.Unauthorized(c, "Unknown token")
	}
}

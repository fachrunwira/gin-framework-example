package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type baseResponse struct {
	Success bool   `json:"success"`
	Msg     string `json:"message"`
}

type successResponse[T any] struct {
	baseResponse
	Content T `json:"contents"`
}

type failedResponse struct {
	baseResponse
	Error string `json:"error"`
}

func Ok[T any](c *gin.Context, msg string, content T) {
	c.JSON(http.StatusOK, successResponse[T]{
		baseResponse: baseResponse{
			Success: true,
			Msg:     msg,
		},
		Content: content,
	})
}

func RouteNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, failedResponse{
		baseResponse: baseResponse{
			Success: false,
			Msg:     "api route is not exists",
		},
		Error: "path_not_found",
	})
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, failedResponse{
		baseResponse: baseResponse{
			Success: false,
			Msg:     msg,
		},
		Error: "unauthorized_access",
	})
}

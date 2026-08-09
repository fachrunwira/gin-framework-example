package appvalidation

import (
	"github.com/gin-gonic/gin"
)

type Validation struct {
	Message string
	Err     error
}

func FormValidation(c *gin.Context, req any, message ErrorMessages) (validation Validation) {
	validation.Err = c.ShouldBind(req)
	if validation.Err != nil {
		validation.Message = CustomFirstError(validation.Err, message)
	}

	return validation
}

func JsonValidation(c *gin.Context, req any, message ErrorMessages) (validation Validation) {
	validation.Err = c.ShouldBindJSON(req)
	if validation.Err != nil {
		validation.Message = CustomFirstError(validation.Err, message)
	}

	return validation
}

func UriValidation(c *gin.Context, req any, message ErrorMessages) (validation Validation) {
	validation.Err = c.ShouldBindUri(req)
	if validation.Err != nil {
		validation.Message = CustomFirstError(validation.Err, message)
	}

	return validation
}

func QueryValidation(c *gin.Context, req any, message ErrorMessages) (validation Validation) {
	validation.Err = c.ShouldBindQuery(req)
	if validation.Err != nil {
		validation.Message = CustomFirstError(validation.Err, message)
	}

	return validation
}
